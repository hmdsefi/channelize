package metrics

import (
	"encoding/hex"
	"math/rand"
	"strings"
	"sync"
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

func TestMetrics(t *testing.T) {
	collector := NewMetrics()

	assert.True(t, strings.Contains(collector.openConnections.Desc().String(), "\"open_connections\""))
}

func TestMetrics_OpenConnection(t *testing.T) {
	t.Run("test open connection", func(t *testing.T) {
		collector := newMetricsWithPostfix(randString())
		assert.Equal(t, float64(0), testutil.ToFloat64(collector.openConnections))
		collector.OpenConnection()
		assert.Equal(t, float64(1), testutil.ToFloat64(collector.openConnections))
	})

	t.Run("test parallel open connection", func(t *testing.T) {
		collector := newMetricsWithPostfix(randString())
		wg := new(sync.WaitGroup)
		n := 10
		wg.Add(n)
		for i := 0; i < n; i++ {
			go func() {
				defer wg.Done()
				collector.OpenConnection()
			}()
		}

		wg.Wait()
		assert.Equal(t, float64(10), testutil.ToFloat64(collector.openConnections))
	})
}

func TestMetrics_CloseConnection(t *testing.T) {
	t.Run("test close connection", func(t *testing.T) {
		collector := newMetricsWithPostfix(randString())
		collector.openConnections.Add(2)
		assert.Equal(t, float64(2), testutil.ToFloat64(collector.openConnections))
		collector.CloseConnection()
		assert.Equal(t, float64(1), testutil.ToFloat64(collector.openConnections))
	})

	t.Run("test parallel close connection", func(t *testing.T) {
		collector := newMetricsWithPostfix(randString())
		wg := new(sync.WaitGroup)
		n := 10
		collector.openConnections.Add(float64(n))
		wg.Add(n)
		for i := 0; i < n; i++ {
			go func() {
				defer wg.Done()
				collector.CloseConnection()
			}()
		}

		wg.Wait()
		assert.Equal(t, float64(0), testutil.ToFloat64(collector.openConnections))
	})
}

func randString() string {
	randBytes := make([]byte, 10)
	rand.Read(randBytes) // nolint
	return hex.EncodeToString(randBytes)
}
