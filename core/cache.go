package core

import (
	"context"
	"sync"

	"github.com/hamed-yousefi/channelize/channel"
	"github.com/hamed-yousefi/channelize/conn"
)

type Cache struct {
	connectionID2Channels map[string]map[channel.Channel]struct{}
	channel2Connections   map[channel.Channel]map[string]*conn.Connection

	sync.RWMutex
}

func NewCache() *Cache {
	return &Cache{
		connectionID2Channels: make(map[string]map[channel.Channel]struct{}),
		channel2Connections:   make(map[channel.Channel]map[string]*conn.Connection),
	}
}

func (c *Cache) Subscribe(_ context.Context, connection *conn.Connection, channels ...channel.Channel) {
	c.Lock()
	defer c.Unlock()

	for _, ch := range channels {
		c.connectionID2Channels[connection.ID()][ch] = struct{}{}
		c.channel2Connections[ch][connection.ID()] = connection
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

func (c *Cache) Connections(ch channel.Channel) []*conn.Connection {
	c.RLock()
	defer c.RUnlock()
	var connections []*conn.Connection
	for _, connection := range c.channel2Connections[ch] {
		connections = append(connections, connection)
	}

	return connections
}
