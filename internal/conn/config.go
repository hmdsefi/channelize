/**
 * Copyright © 2022 Hamed Yousefi <hdyousefi@gmail.com>.
 */

package conn

import (
	"fmt"
	"time"

	"github.com/hmdsefi/channelize/internal/common/utils"
)

const (
	// The default value of buffer size of the outbound channel.
	defaultOutboundBufferSize = 256

	// The default value of reading pong messages.
	defaultPongWait = 60 * time.Second

	// The default value of ping period. It must be less than defaultPongWait.
	defaultPingPeriod = (defaultPongWait * 9) / 10
)

type PingMessageFunc func() []byte

// Config represents the configuration that is needed to create a new Connection.
type Config struct {
	// outboundBufferSize represents the buffer size of the outbound channel.
	outboundBufferSize int

	// pongWait represents the time allowed to read the next pong message from the peer.
	pongWait time.Duration

	// pingPeriod represents the time to send ping to peer. Must be less than pongWait.
	pingPeriod time.Duration

	// pingMessageFunc is a function that create ping messages.
	pingMessageFunc PingMessageFunc

	collector collector
}

func newDefaultConfig() *Config {
	return &Config{
		outboundBufferSize: defaultOutboundBufferSize,
		pongWait:           defaultPongWait,
		pingPeriod:         defaultPingPeriod,
		pingMessageFunc:    defaultPingMessageFunc,
		collector:          newNoopCollector(),
	}
}

type Option func(*Config)

func WithOutboundBufferSize(size int) Option {
	return func(config *Config) {
		if config == nil {
			return
		}

		config.outboundBufferSize = size
	}
}

func WithPongWait(duration time.Duration) Option {
	return func(config *Config) {
		if config == nil {
			return
		}

		config.pongWait = duration
	}
}

func WithPingPeriod(duration time.Duration) Option {
	return func(config *Config) {
		if config == nil {
			return
		}

		config.pingPeriod = duration
	}
}

func WithPingMessageFunc(messageFunc PingMessageFunc) Option {
	return func(config *Config) {
		if config == nil {
			return
		}

		config.pingMessageFunc = messageFunc
	}
}

func WithCollector(in collector) Option {
	return func(config *Config) {
		if config == nil {
			return
		}

		config.collector = in
	}
}

func defaultPingMessageFunc() []byte {
	return []byte(fmt.Sprint(utils.Now().Unix()))
}

type noopCollector struct {
}

func newNoopCollector() *noopCollector {
	return &noopCollector{}
}

func (n *noopCollector) OpenConnectionsInc() {
}

func (n *noopCollector) OpenConnectionsDec() {
}
