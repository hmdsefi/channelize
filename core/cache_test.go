/**
 * Copyright © 2022 Hamed Yousefi <hdyousefi@gmail.com>.
 */

package core

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hamed-yousefi/channelize/channel"
	"github.com/hamed-yousefi/channelize/common"
	"github.com/hamed-yousefi/channelize/core/mock"
)

const (
	testConnID = "test-conn-id"
)

var (
	testChannels = []channel.Channel{
		"error",
		"alerts",
		"notifications",
		"feed",
	}

	testConnectionIDs = []string{
		"test-conn-id-1",
		"test-conn-id-2",
		"test-conn-id-3",
		"test-conn-id-4",
		"test-conn-id-5",
	}
)

// TestCache_Subscribe subscribes a connection to a list of channels.
func TestCache_Subscribe(t *testing.T) {
	ctx := context.Background()
	cache := NewCache()

	expectedConn := mock.NewConnection(testConnID)
	cache.Subscribe(ctx, expectedConn, testChannels[:2]...)

	assert.Equal(t, 1, len(cache.connectionID2Channels))
	assert.Equal(t, 2, len(cache.connectionID2Channels[expectedConn.ID()]))
	require.Equal(t, 2, len(cache.channel2Connections))
	assert.Equal(t, expectedConn, cache.channel2Connections[testChannels[0]][expectedConn.ID()])
	assert.Equal(t, expectedConn, cache.channel2Connections[testChannels[1]][expectedConn.ID()])
}

// TestCache_Unsubscribe unsubscribes from a connection from multiple
// channels concurrently.
func TestCache_Unsubscribe(t *testing.T) {
	ctx := context.Background()
	conn := mock.NewConnection(testConnID)
	cache := initCache(conn)

	for _, ch := range testChannels {
		t.Run("parallel unsubscribe", func(t *testing.T) {
			t.Parallel()
			cache.Unsubscribe(ctx, conn.ID(), ch)
			_, exists := cache.channel2Connections[ch][conn.ID()]
			assert.False(t, exists)

			_, exists = cache.connectionID2Channels[conn.ID()][ch]
			assert.False(t, exists)
		})
	}
}

// TestCache_Remove removes multiple connections from the storage concurrently.
func TestCache_Remove(t *testing.T) {
	ctx := context.Background()
	var connections []common.ConnectionWrapper
	for _, id := range testConnectionIDs {
		connections = append(connections, mock.NewConnection(id))
	}

	cache := initCache(connections...)

	for i := range connections {
		t.Run("parallel remove", func(t *testing.T) {
			t.Parallel()
			cache.Remove(ctx, connections[i].ID())
			_, exists := cache.connectionID2Channels[connections[i].ID()]
			assert.False(t, exists)

			for _, ch := range testChannels {
				_, exists := cache.channel2Connections[ch][connections[i].ID()]
				assert.False(t, exists)
			}
		})
	}
}

// TestCache_Connections returns multiple list of available connections
// for a channel concurrently.
func TestCache_Connections(t *testing.T) {
	ctx := context.Background()
	expectedConnections := map[string]common.ConnectionWrapper{}
	var connections []common.ConnectionWrapper
	for _, id := range testConnectionIDs {
		expectedConnections[id] = mock.NewConnection(id)
		connections = append(connections, expectedConnections[id])
	}

	cache := initCache(connections...)

	for _, ch := range testChannels {
		t.Run("parallel get connections", func(t *testing.T) {
			t.Parallel()
			connections := cache.Connections(ctx, ch)
			actualConnections := map[string]common.ConnectionWrapper{}
			for i := range connections {
				actualConnections[connections[i].ID()] = connections[i]
			}

			for id := range expectedConnections {
				assert.Equal(t, expectedConnections[id], actualConnections[id])
			}
		})
	}
}

func initCache(connections ...common.ConnectionWrapper) *Cache {
	cache := NewCache()
	for i := range connections {
		cache.connectionID2Channels[connections[i].ID()] = make(map[channel.Channel]struct{})

		for _, ch := range testChannels {
			if _, exists := cache.channel2Connections[ch]; !exists {
				cache.channel2Connections[ch] = make(map[string]common.ConnectionWrapper)
			}

			cache.connectionID2Channels[connections[i].ID()][ch] = struct{}{}
			cache.channel2Connections[ch][connections[i].ID()] = connections[i]
		}
	}

	return cache
}
