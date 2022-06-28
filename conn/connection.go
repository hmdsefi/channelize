package conn

import (
	"context"
	"sync"

	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

// Connection wraps the websocket connection and add more functionalities to it.
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
}

func NewConnection(ctx context.Context, conn *websocket.Conn, options ...Option) *Connection {
	// setup connection configuration
	config := newDefaultConfig()

	for _, option := range options {
		option(config)
	}

	// wrap the application context with cancellation
	ctx, cancel := context.WithCancel(ctx)

	out := &Connection{
		id:        uuid.NewV4().String(),
		conn:      conn,
		connected: true,
		cancel:    cancel,
		send:      make(chan []byte, config.outboundBufferSize),
	}

	return out
}

// isConnected returns true if connections. Otherwise, returns false.
func (c *Connection) isConnected() bool {
	c.mu.RLock()
	out := c.connected
	c.mu.RUnlock()

	return out
}

func (c *Connection) Close() error {
	var err error

	// Do it once
	c.once.Do(func() {
		// TODO unsubscribe all the streams

		// cancel the context to exist from read and write goroutines
		c.cancel()

		// close the websocket connection
		err = c.conn.Close()

		// release the memory
		c.conn = nil

		c.mu.Lock()
		c.connected = false
		c.mu.Unlock()
	})

	if err != nil {
		return err
	}

	return nil
}
