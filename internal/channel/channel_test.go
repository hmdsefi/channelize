/**
 * Copyright © 2022 Hamed Yousefi <hdyousefi@gmail.com>.
 */

package channel

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testChannels = []string{
		"news",
		"btcusdt@full",
		"ethusdt@full",
		"xrpusdt@full",
		"trxusdt@full",
		"solusdt@full",
		"btcusdt@miniticker",
		"ethusdt@miniticker",
		"xrpusdt@miniticker",
		"trxusdt@miniticker",
		"solusdt@miniticker",
		"aggregate@miniticker",
	}
)

// TestRegisterPublicChannel registers test channels and check if
// they registered successfully.
func TestRegisterPublicChannel(t *testing.T) {
	var channels []Channel
	for _, channelStr := range testChannels {
		channels = append(channels, RegisterPublicChannel(channelStr))
	}

	for _, channel := range channels {
		assert.True(t, channel.IsSupportedChannel())
		assert.True(t, channel.IsSupportedPublicChannel())
	}
}

// TestRegisterPublicChannels registers test channels and check if
// they registered successfully.
func TestRegisterPublicChannels(t *testing.T) {
	channels := RegisterPublicChannels(testChannels...)

	for _, channel := range channels {
		assert.True(t, channel.IsSupportedChannel())
		assert.True(t, channel.IsSupportedPublicChannel())
	}
}

// TestRegisterPublicChannel_Concurrent registers test channels and start registering
// them again in multiple goroutines and at the same time check if the channels are
// registered or not.
func TestRegisterPublicChannel_Concurrent(t *testing.T) {
	// register all the test channels to make sure they exist
	// during check if they exist or not.
	channels := map[Channel]struct{}{}
	for _, channelStr := range testChannels {
		channels[RegisterPublicChannel(channelStr)] = struct{}{}
	}

	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}
	wg.Add(10)

	// create 10 goroutines that changes the existing channels
	// concurrently in an infinite loop.
	for i := 0; i < 10; i++ {
		idx := i
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					t.Log("Goroutine", idx, "Closed")
					return
				default:
				}
				for _, channelStr := range testChannels {
					RegisterPublicChannel(channelStr)
				}
			}
		}()
	}

	// check if channels are registered when 10 goroutines are running in background.
	for channel := range channels {
		assert.True(t, channel.IsSupportedChannel())
		assert.True(t, channel.IsSupportedPublicChannel())
	}

	// cancel the context and wait for the goroutines to be closed.
	cancel()
	wg.Wait()
}

// TestRegisterPublicChannels_Concurrent registers test channels and start registering
// them again in multiple goroutines and at the same time check if the channels are
// registered or not.
func TestRegisterPublicChannels_Concurrent(t *testing.T) {
	// register all the test channels to make sure they exist
	// during check if they exist or not.
	channels := RegisterPublicChannels(testChannels...)

	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}
	wg.Add(10)

	// create 10 goroutines that changes the existing channels
	// concurrently in an infinite loop.
	for i := 0; i < 10; i++ {
		idx := i
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					t.Log("Goroutine", idx, "Closed")
					return
				default:
				}
				RegisterPublicChannels(testChannels...)
			}
		}()
	}

	// check if channels are registered when 10 goroutines are running in background.
	for i := range channels {
		assert.True(t, channels[i].IsSupportedChannel())
		assert.True(t, channels[i].IsSupportedPublicChannel())
	}

	// cancel the context and wait for the goroutines to be closed.
	cancel()
	wg.Wait()
}
