package conn

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWithOutboundBufferSize(t *testing.T) {
	expectedBufferSize := 11
	option := WithOutboundBufferSize(expectedBufferSize)
	option(nil)

	cfg := newDefaultConfig()
	option(cfg)

	assert.Equal(t, expectedBufferSize, cfg.outboundBufferSize)
}

func TestWithPingPeriod(t *testing.T) {
	expectedPingPeriod := time.Minute
	option := WithPingPeriod(expectedPingPeriod)
	option(nil)

	cfg := newDefaultConfig()
	option(cfg)

	assert.Equal(t, expectedPingPeriod, cfg.pingPeriod)
}

func TestWithPongWait(t *testing.T) {
	expectedPongWait := time.Minute
	option := WithPongWait(expectedPongWait)
	option(nil)

	cfg := newDefaultConfig()
	option(cfg)

	assert.Equal(t, expectedPongWait, cfg.pongWait)
}

func TestWithPingMessageFunc(t *testing.T) {
	expectedPingMessage := []byte("test-ping-message")
	pingFunc := func() []byte {
		return expectedPingMessage
	}
	option := WithPingMessageFunc(pingFunc)
	option(nil)

	cfg := newDefaultConfig()
	option(cfg)

	assert.Equal(t, string(expectedPingMessage), string(cfg.pingMessageFunc()))
}
