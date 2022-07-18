package conn

import (
	"sync/atomic"
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

func TestWithCollector(t *testing.T) {
	c := newMockCollector()
	option := WithCollector(c)
	option(nil)

	cfg := newDefaultConfig()
	option(cfg)

	assert.NotNil(t, cfg.collector)
	cfg.collector.OpenConnection()
	assert.Equal(t, int32(1), c.openConnections)
	cfg.collector.CloseConnection()
	assert.Equal(t, int32(0), c.openConnections)
}

type mockCollector struct {
	openConnections int32
}

func newMockCollector() *mockCollector {
	return &mockCollector{}
}

func (n *mockCollector) OpenConnection() {
	atomic.AddInt32(&n.openConnections, 1)
}

func (n *mockCollector) CloseConnection() {
	atomic.AddInt32(&n.openConnections, -1)
}
