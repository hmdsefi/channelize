package channelize

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"

	"github.com/hamed-yousefi/channelize/internal/channel"
	"github.com/hamed-yousefi/channelize/internal/conn"
	"github.com/hamed-yousefi/channelize/internal/core"
)

// connectionHelper is a middleware between connection and storage. It helps
// to connect a connection to the storage.
//
// connectionHelper gives this ability to the connection to register or unregister
// itself into the storage.
type connectionHelper interface {
	// ParseMessage deserializes the inbound messages and call the storage methods
	// based on the message type.
	ParseMessage(ctx context.Context, connection *conn.Connection, data []byte)

	// Remove removes the connection from the storage.
	Remove(ctx context.Context, connID string)
}

// dispatcher is a mechanism to send the public messages to the existing connections.
type dispatcher interface {
	// SendPublicMessage sends the input message to the connections that already
	// subscribed to the input channel. It doesn't authenticate the connection,
	// since the message is public.
	SendPublicMessage(ctx context.Context, ch channel.Channel, message interface{}) error
}

// Channelize wraps all the internal implementations and restricts the exposed
// functionalities to reduce the public API surface.
//
// It provides more APIs like HTTP handlers to facilitate the API usage.
type Channelize struct {
	helper     connectionHelper
	dispatcher dispatcher
}

// NewChannelize creates new instance of Channelize struct. It uses in-memory
// storage by default to store the connections and mapping between the connections and
// channels.
func NewChannelize() *Channelize {
	storage := core.NewCache()
	return &Channelize{
		helper:     newHelper(storage),
		dispatcher: core.NewDispatch(storage),
	}
}

// CreateConnection creates a `conn.Connection` object with the input options.
func (c *Channelize) CreateConnection(ctx context.Context, wsConn *websocket.Conn, options ...conn.Option) *conn.Connection {
	return conn.NewConnection(ctx, wsConn, c.helper, options...)
}

// MakeEchoHTTPHandler makes an echo HTTP handler function. The client should
// provide the websocket.Upgrader. It automatically creates the websocket.Conn
// and conn.Connection.
func (c *Channelize) MakeEchoHTTPHandler(appCtx context.Context, upgrader websocket.Upgrader, options ...conn.Option) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		wsConn, err := upgrader.Upgrade(ctx.Response().Writer, ctx.Request(), nil)
		if err != nil {
			return err
		}

		c.CreateConnection(appCtx, wsConn, options...)

		return nil
	}
}

// MakeEchoHTTPHandler makes a built-in HTTP handler function. The client should
// provide the websocket.Upgrader. It automatically creates the websocket.Conn
// and conn.Connection.
func (c *Channelize) makeHTTPHandler(appCtx context.Context, upgrader websocket.Upgrader, options ...conn.Option) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wsConn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			// TODO log the error
			http.Error(
				w,
				fmt.Sprintf("failt to create websocket connection: %s", err),
				http.StatusInternalServerError,
			)
			return
		}

		c.CreateConnection(appCtx, wsConn, options...)
	}
}

// RegisterPublicChannel creates and registers a new channel by calling the
// internal channel.RegisterPublicChannel function. It returns the created
// channel.
func RegisterPublicChannel(channelStr string) channel.Channel {
	return channel.RegisterPublicChannel(channelStr)
}

// RegisterPublicChannels creates and registers a list of input channels by
// calling the internal channel.RegisterPublicChannels function. It returns
// a list of created channels.
func RegisterPublicChannels(channels ...string) []channel.Channel {
	return channel.RegisterPublicChannels(channels...)
}
