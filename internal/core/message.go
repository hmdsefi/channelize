/**
 * Copyright © 2022 Hamed Yousefi <hdyousefi@gmail.com>.
 */

package core

import (
	"encoding/json"
	"strings"

	"github.com/hamed-yousefi/channelize/internal/channel"
	"github.com/hamed-yousefi/channelize/internal/common/errorx"
	"github.com/hamed-yousefi/channelize/internal/common/validation"
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
	Token    *string           `json:"token"`
}

// HasToken returns true if token field is not nil or empty string.
func (p paramIn) HasToken() bool {
	return p.Token != nil && len(strings.TrimSpace(*p.Token)) > 0
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
			// check if the channel is supported
			if !ch.IsSupportedChannel() {
				out.AddFieldError(
					validation.SubField(validation.FieldChannels, ch.String()),
					errorx.ErrorMsgUnsupportedChannel,
				)
				continue
			}

			// check if the channel is not private then it should be public
			if !ch.IsSupportedPrivateChannel() && !ch.IsSupportedPublicChannel() {
				out.AddFieldError(
					validation.SubField(validation.FieldChannels, ch.String()),
					errorx.ErrorMsgInvalidChannelType,
				)
			}

			// check if the channel is private, token should exist
			if ch.IsSupportedPrivateChannel() && !m.Params.HasToken() {
				out.AddFieldError(
					validation.SubField(validation.FieldChannels, ch.String()),
					errorx.ErrorMsgAuthTokenIsMissing,
				)
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

func newMessageOut(channel channel.Channel, data interface{}) *MessageOut {
	return &MessageOut{Channel: channel, Data: data}
}
