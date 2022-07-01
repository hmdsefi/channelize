/**
 * Copyright © 2022 Hamed Yousefi <hdyousefi@gmail.com>.
 */

package conn

import (
	"context"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"

	"github.com/hamed-yousefi/channelize/errorx"
	"github.com/hamed-yousefi/channelize/utils"
	"github.com/hamed-yousefi/channelize/validation"
)

// messageProcessor is a mechanism to validation and process peer messages.
type messageProcessor interface {
	Validate(message []byte) validation.Result
	ProcessMessage(ctx context.Context, conn *Connection, message []byte)
}

// Connection wraps the websocket connection and add more functionalities to it.
// Each client that connected to the websocket server has a Connection.
type Connection struct {
	// id represents connectionID
	id string

	// conn represents websocket connection. It is the handshake
	// between the client and the server. Server uses conn to send
	// and receive messages from the client.
	conn *websocket.Conn

	// send is a buffered channel for the outbound messages.
	send chan []byte

	// cancel can close the websocket connection and stop listening
	// and sending messages.
	cancel context.CancelFunc

	// mu locks the connection during closing the connection to prevent
	// errors and panics during listening from or writing to a connection.
	mu sync.RWMutex

	// once won't let the connection closes more than once.
	once sync.Once

	// connected represent the connection status. it's true if the connection
	// is open. Otherwise, it is false.
	connected bool

	// config represents connection configuration.
	config Config
}

// NewConnection creates a new instance of Connection that wraps the input
// websocket.Conn.
//
// It runs read and write goroutines to read the client messages and write
// the server messages to the connection.
//
// websocket.Conn can be close by calling the Connection.Close method. It will
// stop the goroutines and release the memory.
//
// Cancelling input context, closes the connection. So, the input context
// must be the application context not the request context.
func NewConnection(
	ctx context.Context,
	conn *websocket.Conn,
	msgProcessor messageProcessor,
	options ...Option,
) *Connection {
	// setup connection configuration
	config := newDefaultConfig()

	for _, option := range options {
		option(config)
	}

	// wrap the application context with cancellation
	ctx, cancel := context.WithCancel(ctx)

	connWrapper := &Connection{
		id:        uuid.NewV4().String(),
		conn:      conn,
		connected: true,
		cancel:    cancel,
		send:      make(chan []byte, config.outboundBufferSize),
		config:    *config,
	}

	go connWrapper.read(ctx, msgProcessor)
	go connWrapper.write(ctx)

	return connWrapper
}

// ID returns the connection id.
func (c *Connection) ID() string {
	return c.id
}

// SendMessage sends the input message to the outbound channel.
// Connection.write method will receive this message and writes it to the
// client.
//
// Before sending the input message, it checks if the connection is still
// open or not. If it is closed, closes the outbound channel and return error.
//
// Returns error if outbound buffer is full.
func (c *Connection) SendMessage(message []byte) error {
	// check if the connection is already closed, close the outbound
	// channel and return error.
	if !c.isConnected() {
		close(c.send)
		return errorx.NewChannelizeError(errorx.CodeConnectionClosed)
	}

	select {
	case c.send <- message:
		return nil
	default:
		// it happens when Config.outboundBufferSize is too small and load on
		// Connection.SendMessage method is too high.
		return errorx.NewChannelizeError(errorx.CodeOutboundBufferIsFull)
	}
}

// isConnected returns true if connections. Otherwise, returns false.
func (c *Connection) isConnected() bool {
	c.mu.RLock()
	out := c.connected
	c.mu.RUnlock()

	return out
}

// Close closes the websocket.Conn and cancel the connection context.
// Cancelling the context causes closing the running read and write
// goroutines. The closing connection process is singleton.
func (c *Connection) Close() error {
	var err error

	// Do it once
	c.once.Do(func() {
		c.mu.Lock()
		c.connected = false
		c.mu.Unlock()

		// NOTE: do not close the outbound channel here. It can cause panic.
		// The

		// TODO unsubscribe all the streams

		// cancel the context to exist from read and write goroutines
		c.cancel()

		// close the websocket connection
		err = c.conn.Close()
	})

	if err != nil {
		return err
	}

	return nil
}

// read sets the pong handler and read deadline and listen to the
// websocket connection to read the client messages.
//
// If the message type is not close, ping, or pong, the read method
// deserializes and validation the client message.
//
// If the message was valid, it will subscribe or unsubscribe to one
// or more channels.
func (c *Connection) read(ctx context.Context, msgProcessor messageProcessor) {
	defer func() {
		// TODO log the close error
		_ = c.Close()
	}()

	// set pong message expiration time.
	if err := c.conn.SetReadDeadline(utils.Now().Add(c.config.pongWait)); err != nil {
		// TODO log the error
		return
	}

	// create and set pong handler
	c.conn.SetPongHandler(func(in string) error {
		// update pong message expiration time
		if err := c.conn.SetReadDeadline(utils.Now().Add(c.config.pongWait)); err != nil {
			// TODO log the error
			return err
		}

		return nil
	})

	// listen to the websocket connection in an infinite for loop to
	// receive client messages. Break the loop and return if any error
	// happened or received close message from the client.
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		messageType, message, err := c.conn.ReadMessage()
		if err != nil {
			// TODO log the read message error
			return
		}

		// check if the message type is a close message then return.
		if isCloseMessageType(messageType) {
			return
		}

		// check if it is pong message, continue.
		if isPingPongMessageType(messageType) {
			continue
		}

		if v := msgProcessor.Validate(message); !v.IsValid() {
			// TODO log invalid message error
			// TODO should we publish the error to the WS
			continue
		}

		msgProcessor.ProcessMessage(ctx, c, message)
	}
}

func isCloseMessageType(code int) bool {
	return code == websocket.CloseMessage ||
		(code >= websocket.CloseNormalClosure && code <= websocket.CloseTLSHandshake)
}

func isPingPongMessageType(code int) bool {
	return code == websocket.PingMessage || code == websocket.PongMessage
}

// write listens to send channel and write the messages to the
// websocket connection.
//
// It writes a ping message based on the Config.pingPeriod. By
// default, the ping message is unix timestamp.
//
// write will break the for loop and return whenever an error
// happens or context has been cancelled.
func (c *Connection) write(ctx context.Context) {
	pingTicker := time.NewTicker(c.config.pingPeriod)

	defer func() {
		pingTicker.Stop()
		// TODO log the close error
		_ = c.Close()
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case <-pingTicker.C:
			c.mu.RLock()
			// return if the connection is already closed.
			if !c.isConnected() {
				return
			}

			// write the ping message to the peer.
			if err := c.conn.WriteMessage(websocket.PingMessage, c.config.pingMessageFunc()); err != nil {
				// TODO log the error
				return
			}
			c.mu.RUnlock()
		case message := <-c.send:
			c.mu.RLock()
			// return if the connection is already closed.
			if !c.isConnected() {
				return
			}

			// write the message to the peer.
			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				// TODO log the error
				return
			}
			c.mu.RUnlock()
		}
	}
}
