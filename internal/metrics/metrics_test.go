/**
 * Copyright © 2022 Hamed Yousefi <hdyousefi@gmail.com>.
 */

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
	assert.True(t, strings.Contains(collector.privateConnections.Desc().String(), "\"private_connections\""))
}

func TestMetrics_OpenConnectionsInc(t *testing.T) {
	t.Run("test open connection inc", func(t *testing.T) {
		collector := newMetricsWithPostfix(randString())
		assert.Equal(t, float64(0), testutil.ToFloat64(collector.openConnections))
		collector.OpenConnectionsInc()
		assert.Equal(t, float64(1), testutil.ToFloat64(collector.openConnections))
	})

	t.Run("test parallel open connection inc", func(t *testing.T) {
		collector := newMetricsWithPostfix(randString())
		wg := new(sync.WaitGroup)
		n := 10
		wg.Add(n)
		for i := 0; i < n; i++ {
			go func() {
				defer wg.Done()
				collector.OpenConnectionsInc()
			}()
		}

		wg.Wait()
		assert.Equal(t, float64(10), testutil.ToFloat64(collector.openConnections))
	})
}

func TestMetrics_OpenConnectionsDec(t *testing.T) {
	t.Run("test open connection dec", func(t *testing.T) {
		collector := newMetricsWithPostfix(randString())
		collector.openConnections.Add(2)
		assert.Equal(t, float64(2), testutil.ToFloat64(collector.openConnections))
		collector.OpenConnectionsDec()
		assert.Equal(t, float64(1), testutil.ToFloat64(collector.openConnections))
	})

	t.Run("test parallel open connection dec", func(t *testing.T) {
		collector := newMetricsWithPostfix(randString())
		wg := new(sync.WaitGroup)
		n := 10
		collector.openConnections.Add(float64(n))
		wg.Add(n)
		for i := 0; i < n; i++ {
			go func() {
				defer wg.Done()
				collector.OpenConnectionsDec()
			}()
		}

		wg.Wait()
		assert.Equal(t, float64(0), testutil.ToFloat64(collector.openConnections))
	})
}

func TestMetrics_PrivateConnectionsInc(t *testing.T) {
	t.Run("test private connection inc", func(t *testing.T) {
		collector := newMetricsWithPostfix(randString())
		assert.Equal(t, float64(0), testutil.ToFloat64(collector.privateConnections))
		collector.PrivateConnectionsInc()
		assert.Equal(t, float64(1), testutil.ToFloat64(collector.privateConnections))
	})

	t.Run("test parallel private connection inc", func(t *testing.T) {
		collector := newMetricsWithPostfix(randString())
		wg := new(sync.WaitGroup)
		n := 10
		wg.Add(n)
		for i := 0; i < n; i++ {
			go func() {
				defer wg.Done()
				collector.PrivateConnectionsInc()
			}()
		}

		wg.Wait()
		assert.Equal(t, float64(10), testutil.ToFloat64(collector.privateConnections))
	})
}

func TestMetrics_PrivateConnectionsDec(t *testing.T) {
	t.Run("test private connection dec", func(t *testing.T) {
		collector := newMetricsWithPostfix(randString())
		collector.privateConnections.Add(2)
		assert.Equal(t, float64(2), testutil.ToFloat64(collector.privateConnections))
		collector.PrivateConnectionsDec()
		assert.Equal(t, float64(1), testutil.ToFloat64(collector.privateConnections))
	})

	t.Run("test parallel private connection dec", func(t *testing.T) {
		collector := newMetricsWithPostfix(randString())
		wg := new(sync.WaitGroup)
		n := 10
		collector.privateConnections.Add(float64(n))
		wg.Add(n)
		for i := 0; i < n; i++ {
			go func() {
				defer wg.Done()
				collector.PrivateConnectionsDec()
			}()
		}

		wg.Wait()
		assert.Equal(t, float64(0), testutil.ToFloat64(collector.privateConnections))
	})
}

func randString() string {
	randBytes := make([]byte, 10)
	rand.Read(randBytes) // nolint
	return hex.EncodeToString(randBytes)
}
