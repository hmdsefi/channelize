package core

import (
	"context"
	"sync"

	"github.com/hamed-yousefi/channelize/channel"
	"github.com/hamed-yousefi/channelize/common"
)

type Cache struct {
	connectionID2Channels map[string]map[channel.Channel]struct{}
	channel2Connections   map[channel.Channel]map[string]common.ConnectionWrapper

	sync.RWMutex
}

func NewCache() *Cache {
	return &Cache{
		connectionID2Channels: make(map[string]map[channel.Channel]struct{}),
		channel2Connections:   make(map[channel.Channel]map[string]common.ConnectionWrapper),
	}
}

func (c *Cache) Subscribe(_ context.Context, conn common.ConnectionWrapper, channels ...channel.Channel) {
	c.Lock()
	defer c.Unlock()

	if _, exists := c.connectionID2Channels[conn.ID()]; !exists {
		c.connectionID2Channels[conn.ID()] = make(map[channel.Channel]struct{})
	}

	for _, ch := range channels {
		if _, exists := c.channel2Connections[ch]; !exists {
			c.channel2Connections[ch] = make(map[string]common.ConnectionWrapper)
		}

		c.connectionID2Channels[conn.ID()][ch] = struct{}{}
		c.channel2Connections[ch][conn.ID()] = conn
	}
}

func (c *Cache) Unsubscribe(_ context.Context, connID string, channels ...channel.Channel) {
	c.Lock()
	defer c.Unlock()

	for _, ch := range channels {
		delete(c.connectionID2Channels[connID], ch)
		delete(c.channel2Connections[ch], connID)
	}
}

func (c *Cache) Remove(_ context.Context, connID string) {
	c.Lock()
	defer c.Unlock()

	for ch := range c.connectionID2Channels[connID] {
		delete(c.channel2Connections[ch], connID)
	}

	delete(c.connectionID2Channels, connID)
}

func (c *Cache) Connections(ch channel.Channel) []common.ConnectionWrapper {
	c.RLock()
	defer c.RUnlock()
	var connections []common.ConnectionWrapper
	for _, conn := range c.channel2Connections[ch] {
		connections = append(connections, conn)
	}

	return connections
}
