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

	supportedChannels        = map[Channel]struct{}{}
	supportedPublicChannels  = map[Channel]struct{}{}
	supportedPrivateChannels = map[Channel]struct{}{}
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

// IsSupportedPublicChannel checks if the channel value is a valid public
// channel or not. It is trade-safe.
func (c Channel) IsSupportedPublicChannel() bool {
	mu.RLock()
	defer mu.RUnlock()
	_, ok := supportedPublicChannels[c]
	return ok
}

// IsSupportedPrivateChannel checks if the channel value is a valid private
// channel or not. It is trade-safe.
func (c Channel) IsSupportedPrivateChannel() bool {
	mu.RLock()
	defer mu.RUnlock()
	_, ok := supportedPrivateChannels[c]
	return ok
}

// RegisterPublicChannel registers a new public channel. It converts the input string
// to the Channel type and adds it to the supportedChannels a supportedPublicChannels
// maps.
//
// RegisterPublicChannel is thread safe and client can use it in multiple goroutines.
//
// Client should call this function in application startup to register the public channels.
func RegisterPublicChannel(channelStr string) Channel {
	mu.Lock()
	defer mu.Unlock()

	channel := Channel(channelStr)
	supportedChannels[channel] = struct{}{}
	supportedPublicChannels[channel] = struct{}{}

	return channel
}

// RegisterPublicChannels registers a list of public channels. It is thread safe.
func RegisterPublicChannels(channels ...string) []Channel {
	mu.Lock()
	defer mu.Unlock()

	out := make([]Channel, len(channels))
	for i := range channels {
		out[i] = Channel(channels[i])
		supportedChannels[out[i]] = struct{}{}
		supportedPublicChannels[out[i]] = struct{}{}
	}

	return out
}

// RegisterPrivateChannel registers a new private channel. It converts the input string
// to the Channel type and adds it to the supportedChannels a supportedPrivateChannels
// maps.
//
// RegisterPrivateChannel is thread safe and client can use it in multiple goroutines.
//
// Client should call this function in application startup to register the private channels.
func RegisterPrivateChannel(channelStr string) Channel {
	mu.Lock()
	defer mu.Unlock()

	channel := Channel(channelStr)
	supportedChannels[channel] = struct{}{}
	supportedPrivateChannels[channel] = struct{}{}

	return channel
}

// RegisterPrivateChannels registers a list of private channels. It is thread safe.
func RegisterPrivateChannels(channels ...string) []Channel {
	mu.Lock()
	defer mu.Unlock()

	out := make([]Channel, len(channels))
	for i := range channels {
		out[i] = Channel(channels[i])
		supportedChannels[out[i]] = struct{}{}
		supportedPrivateChannels[out[i]] = struct{}{}
	}

	return out
}
