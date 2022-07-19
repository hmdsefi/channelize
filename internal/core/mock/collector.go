/**
 * Copyright © 2022 Hamed Yousefi <hdyousefi@gmail.com>.
 */

package mock

import "sync/atomic"

type Collector struct {
	privateConnections int32
}

func NewCollector() *Collector {
	return &Collector{}
}

func (c *Collector) PrivateConnectionsInc() {
	atomic.AddInt32(&c.privateConnections, 1)
}

func (c *Collector) PrivateConnectionsDec() {
	atomic.AddInt32(&c.privateConnections, -1)
}

func (c Collector) PrivateConnections() int32 {
	return c.privateConnections
}
