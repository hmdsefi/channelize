package core

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hamed-yousefi/channelize/channel"
	"github.com/hamed-yousefi/channelize/core/mock"
)

const (
	testConnID = "test-conn-id"
)

var (
	testChannels = []string{
		"error",
		"alerts",
		"notifications",
		"feed",
	}
)

func TestCache_Subscribe(t *testing.T) {
	ctx := context.Background()
	cache := NewCache()
	channels := channel.RegisterPublicChannels(testChannels...)

	expectedConn := mock.NewConnection(testConnID)
	cache.Subscribe(ctx, expectedConn, channels[:2]...)

	assert.Equal(t, 1, len(cache.connectionID2Channels))
	assert.Equal(t, 2, len(cache.connectionID2Channels[expectedConn.ID()]))
	require.Equal(t, 2, len(cache.channel2Connections))
	assert.Equal(t, expectedConn, cache.channel2Connections[channels[0]][expectedConn.ID()])
	assert.Equal(t, expectedConn, cache.channel2Connections[channels[1]][expectedConn.ID()])
}
