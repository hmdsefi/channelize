/**
 * Copyright © 2022 Hamed Yousefi <hdyousefi@gmail.com>.
 */

package core

import (
	"context"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hmdsefi/channelize/internal/channel"
	"github.com/hmdsefi/channelize/internal/common"
	"github.com/hmdsefi/channelize/internal/core/mock"
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

func authNoopFunc() error {
	return nil
}

// TestCache_Subscribe subscribes a connection to a list of channels.
func TestCache_Subscribe(t *testing.T) {
	ctx := context.Background()

	t.Run("subscribe public channels", func(t *testing.T) {
		cache := NewCache()
		expectedConn := mock.NewConnection(testConnID, nil, authNoopFunc)
		cache.Subscribe(ctx, expectedConn, testChannels[:2]...)

		assert.Equal(t, 1, len(cache.connectionID2Channels))
		assert.Equal(t, 2, len(cache.connectionID2Channels[expectedConn.ID()]))
		require.Equal(t, 2, len(cache.channel2Connections))
		assert.Equal(t, expectedConn, cache.channel2Connections[testChannels[0]][expectedConn.ID()])
		assert.Equal(t, expectedConn, cache.channel2Connections[testChannels[1]][expectedConn.ID()])
		assert.True(t, len(cache.userID2ConnectionID) == 0)
	})

	t.Run("subscribe private channels", func(t *testing.T) {
		cache := NewCache()
		userID := uuid.NewV4().String()
		expectedConn := mock.NewConnection(testConnID, &userID, authNoopFunc)
		cache.Subscribe(ctx, expectedConn, testChannels[2:]...)

		assert.Equal(t, 1, len(cache.connectionID2Channels))
		assert.Equal(t, 2, len(cache.connectionID2Channels[expectedConn.ID()]))
		require.Equal(t, 2, len(cache.channel2Connections))
		assert.Equal(t, expectedConn, cache.channel2Connections[testChannels[2]][expectedConn.ID()])
		assert.Equal(t, expectedConn, cache.channel2Connections[testChannels[3]][expectedConn.ID()])
		assert.Equal(t, expectedConn.ID(), cache.userID2ConnectionID[userID])
	})
}

// TestCache_Unsubscribe unsubscribes from a connection from multiple
// channels concurrently.
func TestCache_Unsubscribe(t *testing.T) {
	ctx := context.Background()
	conn := mock.NewConnection(testConnID, nil, authNoopFunc)
	cache := initCache(conn)

	for _, ch := range testChannels {
		t.Run("parallel unsubscribe", func(t *testing.T) {
			t.Parallel()
			cache.Unsubscribe(ctx, conn.ID(), ch)

			cache.RLock()
			defer cache.RUnlock()

			_, exists := cache.channel2Connections[ch][conn.ID()]
			assert.False(t, exists)

			_, exists = cache.connectionID2Channels[conn.ID()][ch]
			assert.False(t, exists)
		})
	}
}

// TestCache_UnsubscribeUserID unsubscribes userIDs from multiple channels concurrently.
func TestCache_UnsubscribeUserID(t *testing.T) {
	ctx := context.Background()
	var connections []common.ConnectionWrapper
	userID2Connection := map[string]common.ConnectionWrapper{}
	for _, id := range testConnectionIDs {
		userID := uuid.NewV4().String()
		conn := mock.NewConnection(id, &userID, authNoopFunc)
		connections = append(connections, conn)
		userID2Connection[userID] = conn
	}

	cache := initCache(connections...)

	for _, ch := range testChannels {
		for userID := range userID2Connection {
			userID := userID
			t.Run("parallel unsubscribe by userID", func(t *testing.T) {
				t.Parallel()
				cache.UnsubscribeUserID(ctx, userID2Connection[userID].ID(), userID, ch)

				cache.RLock()
				defer cache.RUnlock()

				_, exists := cache.channel2Connections[ch][userID2Connection[userID].ID()]
				assert.False(t, exists)

				_, exists = cache.connectionID2Channels[userID2Connection[userID].ID()][ch]
				assert.False(t, exists)

				_, exists = cache.userID2ConnectionID[userID]
				assert.False(t, exists)
			})
		}
	}
}

// TestCache_Remove removes multiple connections from the storage concurrently.
func TestCache_Remove(t *testing.T) {
	ctx := context.Background()
	var connections []common.ConnectionWrapper
	for _, id := range testConnectionIDs {
		userID := uuid.NewV4().String()
		connections = append(connections, mock.NewConnection(id, &userID, authNoopFunc))
	}

	cache := initCache(connections...)

	for i := range connections {
		t.Run("parallel remove", func(t *testing.T) {
			t.Parallel()
			cache.Remove(ctx, connections[i].ID(), connections[i].UserID())

			cache.RLock()
			defer cache.RUnlock()

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
	var testConnections []common.ConnectionWrapper
	for _, id := range testConnectionIDs {
		expectedConnections[id] = mock.NewConnection(id, nil, authNoopFunc)
		testConnections = append(testConnections, expectedConnections[id])
	}

	cache := initCache(testConnections...)

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

		userID := connections[i].UserID()
		if userID != nil {
			cache.userID2ConnectionID[*userID] = connections[i].ID()
		}

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

// TestCache_ConnectionByUserID returns multiple list of available connections
// for a channel concurrently.
func TestCache_ConnectionByUserID(t *testing.T) {
	ctx := context.Background()
	var connections []common.ConnectionWrapper
	userID2Connection := map[string]common.ConnectionWrapper{}
	for _, id := range testConnectionIDs {
		userID := uuid.NewV4().String()
		connection := mock.NewConnection(id, &userID, authNoopFunc)
		connections = append(connections, connection)
		userID2Connection[userID] = connection
	}

	cache := initCache(connections...)

	for _, ch := range testChannels {
		for userID := range userID2Connection {
			expectedUserID := userID
			t.Run("parallel get connections", func(t *testing.T) {
				t.Parallel()
				actualConn := cache.ConnectionByUserID(ctx, ch, expectedUserID)
				actualUserID := actualConn.UserID()
				require.NotNil(t, actualUserID)
				assert.Equal(t, expectedUserID, *actualUserID)
				assert.Equal(t, userID2Connection[expectedUserID].ID(), actualConn.ID())
			})
		}
	}
}
