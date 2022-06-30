/**
 * Copyright © 2022 Hamed Yousefi <hdyousefi@gmail.com>.
 */

package channel

import "sync"

const (
	// ErrorChannel handles all the errors that happens inside the server.
	ErrorChannel Channel = "error"
)

// Channel represents a websocket stream channel
type Channel string

var (
	mu = sync.RWMutex{}

	supportedChannels       = map[Channel]struct{}{}
	supportedPublicChannels = map[Channel]struct{}{}
)

func (c Channel) String() string {
	return string(c)
}

// IsSupportedChannel checks if the channel value is valid or not.
// It is trade-safe.
func (c Channel) IsSupportedChannel() bool {
	mu.RLock()
	defer mu.RUnlock()
	_, ok := supportedChannels[c]
	return ok
}

// IsSupportedPublicChannel checks if the channel value is valid public
// channel or not. It is trade-safe.
func (c Channel) IsSupportedPublicChannel() bool {
	mu.RLock()
	defer mu.RUnlock()
	_, ok := supportedPublicChannels[c]
	return ok
}

// RegisterPublicChannel registers a new channel. It converts the input stream
// to the Channel type and adds it to the supportedChannels a supportedPublicChannels
// maps.
//
// RegisterPublicChannel is thread safe and client can use it in multiple goroutines.
//
// Client should call this function in application startup to register the channels.
func RegisterPublicChannel(channelStr string) Channel {
	mu.Lock()
	defer mu.Unlock()

	channel := Channel(channelStr)
	supportedChannels[channel] = struct{}{}
	supportedPublicChannels[channel] = struct{}{}

	return channel
}
