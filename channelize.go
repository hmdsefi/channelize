/**
 * Copyright © 2022 Hamed Yousefi <hdyousefi@gmail.com>.
 */

package channelize

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"

	"github.com/hamed-yousefi/channelize/internal/channel"
	"github.com/hamed-yousefi/channelize/internal/common"
	internalLog "github.com/hamed-yousefi/channelize/internal/common/log"
	"github.com/hamed-yousefi/channelize/internal/conn"
	"github.com/hamed-yousefi/channelize/internal/core"
	"github.com/hamed-yousefi/channelize/log"
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

type Option func(*Config)

type Config struct {
	logger log.Logger
}

func newDefaultConfig() *Config {
	return &Config{
		logger: internalLog.NewDefaultLogger(),
	}
}

func WithLogger(logger log.Logger) func(config *Config) {
	return func(config *Config) {
		config.logger = logger
	}
}

// Channelize wraps all the internal implementations and restricts the exposed
// functionalities to reduce the public API surface.
//
// It provides more APIs like HTTP handlers to facilitate the API usage.
type Channelize struct {
	helper     connectionHelper
	dispatcher dispatcher
	logger     log.Logger
}

// NewChannelize creates new instance of Channelize struct. It uses in-memory
// storage by default to store the connections and mapping between the connections and
// channels.
func NewChannelize(options ...Option) *Channelize {
	config := newDefaultConfig()
	for _, option := range options {
		option(config)
	}

	storage := core.NewCache()

	return &Channelize{
		helper:     newHelper(storage),
		dispatcher: core.NewDispatch(storage, config.logger),
		logger:     config.logger,
	}
}

// CreateConnection creates a `conn.Connection` object with the input options.
func (c *Channelize) CreateConnection(ctx context.Context, wsConn *websocket.Conn, options ...conn.Option) *conn.Connection {
	return conn.NewConnection(ctx, wsConn, c.helper, c.logger, options...)
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

// MakeHTTPHandler makes a built-in HTTP handler function. The client should
// provide the websocket.Upgrader. It automatically creates the websocket.Conn
// and conn.Connection.
func (c *Channelize) MakeHTTPHandler(appCtx context.Context, upgrader websocket.Upgrader, options ...conn.Option) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wsConn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			c.logger.Error("failed to create websocket.Conn", common.LogFieldError, err.Error())
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

// WithOutboundBufferSize sets the outbound buffer size.
func WithOutboundBufferSize(size int) conn.Option {
	return conn.WithOutboundBufferSize(size)
}

// WithPongWait sets the pong wait duration.
func WithPongWait(duration time.Duration) conn.Option {
	return conn.WithPongWait(duration)
}

// WithPingPeriod sets the ping period.
func WithPingPeriod(duration time.Duration) conn.Option {
	return conn.WithPingPeriod(duration)
}

// WithPingMessageFunc sets the ping function. Client send customized ping messages.
func WithPingMessageFunc(messageFunc conn.PingMessageFunc) conn.Option {
	return conn.WithPingMessageFunc(messageFunc)
}