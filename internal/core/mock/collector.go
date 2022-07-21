/**
 * Copyright © 2022 Hamed Yousefi <hdyousefi@gmail.com>.
 */

package mock

import (
	"sync"
	"sync/atomic"
)

type atomicFloat64 struct {
	sync.RWMutex
	value float64
}

func (a *atomicFloat64) Set(val float64) {
	a.Lock()
	defer a.Unlock()
	a.value = val
}

func (a *atomicFloat64) Value() float64 {
	a.RLock()
	defer a.RUnlock()
	return a.value
}

type Collector struct {
	PrivateConnectionsGauge int32
	SubscribedChannelsCount *atomicFloat64
	OpenConnectionsCount    *atomicFloat64
	PrivateConnectionsCount *atomicFloat64
}

func NewCollector() *Collector {
	return &Collector{
		SubscribedChannelsCount: new(atomicFloat64),
		OpenConnectionsCount:    new(atomicFloat64),
		PrivateConnectionsCount: new(atomicFloat64),
	}
}

func (c *Collector) PrivateConnectionsInc() {
	atomic.AddInt32(&c.PrivateConnectionsGauge, 1)
}

func (c *Collector) PrivateConnectionsDec() {
	atomic.AddInt32(&c.PrivateConnectionsGauge, -1)
}

func (c *Collector) SubscribedChannels(val float64) {
	c.SubscribedChannelsCount.Set(val)
}

func (c *Collector) PrivateConnections(val float64) {
	c.PrivateConnectionsCount.Set(val)
}

func (c *Collector) OpenConnections(val float64) {
	c.OpenConnectionsCount.Set(val)
}
