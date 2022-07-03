/**
 * Copyright © 2022 Hamed Yousefi <hdyousefi@gmail.com>.
 */

package core

import (
	"context"
	"sync"

	"github.com/hamed-yousefi/channelize/channel"
	"github.com/hamed-yousefi/channelize/common"
)

// Cache is an in-memory storage to store available channels and connections.
type Cache struct {
	// connectionID2Channels stores a mapping between the connection ID and channels.
	// It stores list of channels as a map to increase the performance and guaranty
	// the uniqueness of subscribed channels per connection ID.
	// The key of the first map is connection ID and the key of second map is channel.
	// map[connID]map[channel]struct{}
	connectionID2Channels map[string]map[channel.Channel]struct{}

	// channel2Connections stores a list of connections that subscribed to a channel.
	// It's a reverse storage of connectionID2Channels to facilitate finding the
	// connection per channel.
	// The key of first map is channel and the key of second map is connection ID.
	// map[channel]map[connID]connection
	channel2Connections map[channel.Channel]map[string]common.ConnectionWrapper

	sync.RWMutex
}

// NewCache creates a new instance on Cache.
func NewCache() *Cache {
	return &Cache{
		connectionID2Channels: make(map[string]map[channel.Channel]struct{}),
		channel2Connections:   make(map[channel.Channel]map[string]common.ConnectionWrapper),
	}
}

// Subscribe stores the subscription for the input connection and list
// of the channels into the internal maps.
//
// This function is thread-safe and multiple goroutines can subscribe to
// a list of channel concurrently.
func (c *Cache) Subscribe(_ context.Context, conn common.ConnectionWrapper, channels ...channel.Channel) {
	c.Lock()
	defer c.Unlock()

	// check if connection doesn't subscribe to a channel yet, then create
	// the channel map for it to prevent nil pointer panic.
	if _, exists := c.connectionID2Channels[conn.ID()]; !exists {
		c.connectionID2Channels[conn.ID()] = make(map[channel.Channel]struct{})
	}

	// iterate over the input channel and store the subscription.
	for _, ch := range channels {
		if _, exists := c.channel2Connections[ch]; !exists {
			c.channel2Connections[ch] = make(map[string]common.ConnectionWrapper)
		}

		c.connectionID2Channels[conn.ID()][ch] = struct{}{}
		c.channel2Connections[ch][conn.ID()] = conn
	}
}

// Unsubscribe removes the input channels subscription from the internal maps.
//
// This function is thread-safe and multiple goroutines can unsubscribe
// a list of channel concurrently.
func (c *Cache) Unsubscribe(_ context.Context, connID string, channels ...channel.Channel) {
	c.Lock()
	defer c.Unlock()

	for _, ch := range channels {
		delete(c.connectionID2Channels[connID], ch)
		delete(c.channel2Connections[ch], connID)
	}
}

// Remove removes all subscriptions of the input connection id. Removing
// all subscription means removing connection from the storage.
//
// This function is thread-safe and multiple goroutines can remove
// a connection from the in-memory storage.
func (c *Cache) Remove(_ context.Context, connID string) {
	c.Lock()
	defer c.Unlock()

	for ch := range c.connectionID2Channels[connID] {
		delete(c.channel2Connections[ch], connID)
	}

	delete(c.connectionID2Channels, connID)
}

// Connections returns a list of connections that already subscribed
// to the input channel.
//
// This function is thread-safe and multiple goroutines can get the
// list of subscribed connection concurrently.
func (c *Cache) Connections(_ context.Context, ch channel.Channel) []common.ConnectionWrapper {
	c.RLock()
	defer c.RUnlock()

	var connections []common.ConnectionWrapper
	for _, conn := range c.channel2Connections[ch] {
		connections = append(connections, conn)
	}

	return connections
}
