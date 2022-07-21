/**
 * Copyright © 2022 Hamed Yousefi <hdyousefi@gmail.com>.
 */

package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

// Metrics represents application metrics. It is responsible to manages the
// application metrics and registers them in prometheus.
type Metrics struct {
	// openConnections represents total number of open connections.
	openConnections prometheus.Gauge

	// privateConnections represents total number of private connections.
	privateConnections prometheus.Gauge

	subscribedChannelsSet prometheus.Gauge
	openConnectionsSet    prometheus.Gauge
	privateConnectionsSet prometheus.Gauge
}

func NewMetrics() *Metrics {
	return newMetricsWithPostfix("")
}

func newMetricsWithPostfix(postfix string) *Metrics {
	openConnections := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "open_connections" + postfix,
		Help: "Total number of open connections",
	})

	privateConnections := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "private_connections" + postfix,
		Help: "Total number of private connections",
	})

	privateConnectionsSet := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "private_connections_storage_length" + postfix,
		Help: "Total number of private connections based on the length of storage",
	})

	openConnectionsSet := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "open_connections_storage_length" + postfix,
		Help: "Total number of open connections based on the length of storage",
	})

	subscribedChannelsSet := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "subscribed_channels_storage_length" + postfix,
		Help: "Total number of open connections based on the length of storage",
	})

	prometheus.MustRegister(openConnections)

	return &Metrics{
		openConnections:       openConnections,
		privateConnections:    privateConnections,
		privateConnectionsSet: privateConnectionsSet,
		openConnectionsSet:    openConnectionsSet,
		subscribedChannelsSet: subscribedChannelsSet,
	}
}

// OpenConnectionsInc increases the total number of open connections.
func (m *Metrics) OpenConnectionsInc() {
	m.openConnections.Inc()
}

// OpenConnectionsDec decreases the total number of open connections.
func (m *Metrics) OpenConnectionsDec() {
	m.openConnections.Dec()
}

// PrivateConnectionsInc increases the total number of private connections.
func (m *Metrics) PrivateConnectionsInc() {
	m.privateConnections.Inc()
}

// PrivateConnectionsDec decreases the total number of private connections.
func (m *Metrics) PrivateConnectionsDec() {
	m.privateConnections.Dec()
}

func (m *Metrics) PrivateConnections(in float64) {
	m.privateConnectionsSet.Set(in)
}

func (m *Metrics) OpenConnections(in float64) {
	m.openConnectionsSet.Set(in)
}

func (m *Metrics) SubscribedChannels(in float64) {
	m.subscribedChannelsSet.Set(in)
}
