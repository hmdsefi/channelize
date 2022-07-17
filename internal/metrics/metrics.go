package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	openConnections prometheus.Gauge
}

func NewMetrics() *Metrics {
	openConnections := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "open_connections",
		Help: "Number of open connections",
	})

	prometheus.MustRegister(openConnections)

	return &Metrics{
		openConnections: openConnections,
	}
}

func (m *Metrics) OpenConnection() {
	m.openConnections.Inc()
}

func (m *Metrics) CloseConnection() {
	m.openConnections.Dec()
}
