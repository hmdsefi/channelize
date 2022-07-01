/**
 * Copyright © 2022 Hamed Yousefi <hdyousefi@gmail.com>.
 */

package core

import (
	"encoding/json"

	"github.com/hamed-yousefi/channelize/channel"
	"github.com/hamed-yousefi/channelize/errorx"
	"github.com/hamed-yousefi/channelize/validation"
)

const (
	// MessageTypeSubscribe subscribes client to one or more requested channels.
	MessageTypeSubscribe MessageType = "subscribe"

	// MessageTypeUnsubscribe unsubscribes client from the channels that has been
	// already subscribed by the client.
	MessageTypeUnsubscribe MessageType = "unsubscribe"
)

var (
	supportedMessageTypes = map[MessageType]struct{}{
		MessageTypeSubscribe:   {},
		MessageTypeUnsubscribe: {},
	}
)

// MessageType is an alias type of string that represent client message type.
// It describes the client action.
type MessageType string

func (m MessageType) String() string {
	return string(m)
}

// isSupportedMessageType returns true if the message type is supported.
// Otherwise, returns false.
func (m MessageType) isSupportedMessageType() bool {
	_, ok := supportedMessageTypes[m]
	return ok
}

type paramIn struct {
	Channels []channel.Channel `json:"channels"`
}

// messageIn represents the inbound message. It includes an action and
// some parameters that server needs to do the action.
//
// An action is a MessageType and parameters stored in paramsIn struct.
type messageIn struct {
	MessageType MessageType `json:"type"`
	Params      paramIn     `json:"params"`
}

// UnmarshalMessageIn deserializes the input slice of bytes that has
// been read from the websocket connection.
func UnmarshalMessageIn(data []byte) (*messageIn, error) {
	var msgIn messageIn
	if err := json.Unmarshal(data, &msgIn); err != nil {
		return nil, errorx.NewChannelizeErrorWithErr(errorx.CodeFailedToUnmarshalMessage, err)
	}

	return &msgIn, nil
}

// Validate validates all the fields that client sent to the server.
// Input parameters should be matched with action.
func (m messageIn) Validate() *validation.Result {
	out := new(validation.Result)

	if !m.MessageType.isSupportedMessageType() {
		out.AddFieldError(validation.FieldType, errorx.ErrorMsgUnsupportedMessageType)
	}

	if len(m.Params.Channels) == 0 {
		out.AddFieldError(validation.FieldChannels, errorx.ErrorMsgChannelsIsEmpty)
	}

	if len(m.Params.Channels) != 0 {
		for _, ch := range m.Params.Channels {
			if !ch.IsSupportedChannel() {
				out.AddFieldError(validation.FieldChannels+":"+ch.String(), errorx.ErrorMsgChannelsIsEmpty)
			}
		}
	}

	return out
}

// MessageOut represents the outbound message. Each the outbound message
// includes a channel name that the message belongs to it, and the data
// that is the main content.
type MessageOut struct {
	Channel channel.Channel `json:"channel"`
	Data    interface{}     `json:"data"`
}
