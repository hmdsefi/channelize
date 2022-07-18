package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

// Metrics represents application metrics. It is responsible to manages the
// application metrics and registers them in prometheus.
type Metrics struct {
	// openConnections represents total number of open connections.
	openConnections prometheus.Gauge
}

func NewMetrics() *Metrics {
	return newMetricsWithPostfix("")
}

func newMetricsWithPostfix(postfix string) *Metrics {
	openConnections := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "open_connections" + postfix,
		Help: "Number of open connections",
	})

	prometheus.MustRegister(openConnections)

	return &Metrics{
		openConnections: openConnections,
	}
}

// OpenConnection increases the total number of open connections.
func (m *Metrics) OpenConnection() {
	m.openConnections.Inc()
}

// CloseConnection decreases the total number of open connections.
func (m *Metrics) CloseConnection() {
	m.openConnections.Dec()
}
