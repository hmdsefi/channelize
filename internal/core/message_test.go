/**
 * Copyright © 2022 Hamed Yousefi <hdyousefi@gmail.com>.
 */

package core

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hamed-yousefi/channelize/internal/channel"
	"github.com/hamed-yousefi/channelize/internal/common/errorx"
	"github.com/hamed-yousefi/channelize/internal/common/validation"
)

// TestUnmarshalMessageIn unmarshal a correct message and a message with json syntax error.
func TestUnmarshalMessageIn(t *testing.T) {
	expectedErr := "failed to unmarshal inbound message: invalid character ':' after object key:value pair"
	expectedErrCode := errorx.CodeFailedToUnmarshalMessage

	expectedMsg := &messageIn{
		MessageType: MessageTypeSubscribe,
		Params: paramIn{
			Channels: []channel.Channel{
				"notification",
				"alert",
				"feed",
			},
		},
	}

	correctJSONString := `{"type":"subscribe","params":{"channels":["notification","alert","feed"]}}`
	incorrectJSONString := `{"type":"subscribe","params":"channels":["notification","feed"]}}`

	t.Run("unmarshal correct json input", func(t *testing.T) {
		msg, err := UnmarshalMessageIn([]byte(correctJSONString))
		require.Nil(t, err)
		assert.Equal(t, expectedMsg, msg)
	})

	t.Run("unmarshal incorrect json input", func(t *testing.T) {
		msg, err := UnmarshalMessageIn([]byte(incorrectJSONString))
		assert.NotNil(t, err)
		var chanErr *errorx.ChannelizeError
		assert.True(t, errors.As(err, &chanErr))
		assert.Equal(t, expectedErr, err.Error())
		assert.Equal(t, expectedErrCode, chanErr.Code)
		assert.Nil(t, msg)
	})
}

// TestMessageIn_Validate registers a set of channels and test validation of different messages.
func TestMessageIn_Validate(t *testing.T) {
	channels := registerChannels()

	t.Run("valid messageIn", func(t *testing.T) {
		invalidMsg := messageIn{
			MessageType: MessageTypeSubscribe,
			Params: paramIn{
				Channels: channels,
			},
		}

		result := invalidMsg.Validate()
		expectedResult := new(validation.Result)
		assert.Equal(t, expectedResult, result)
	})

	t.Run("invalid messageIn: missing channels", func(t *testing.T) {
		invalidMsg := messageIn{
			MessageType: MessageTypeSubscribe,
		}

		result := invalidMsg.Validate()
		expectedResult := new(validation.Result)
		expectedResult.AddFieldError(validation.FieldChannels, errorx.ErrorMsgChannelsIsEmpty)
		assert.Equal(t, expectedResult, result)
	})

	t.Run("invalid messageIn: unsupported message type", func(t *testing.T) {
		invalidMsg := messageIn{
			MessageType: MessageType("my-type"),
			Params: paramIn{
				Channels: channels,
			},
		}

		result := invalidMsg.Validate()
		expectedResult := new(validation.Result)
		expectedResult.AddFieldError(validation.FieldType, errorx.ErrorMsgUnsupportedMessageType)
		assert.Equal(t, expectedResult, result)
	})

	t.Run("invalid messageIn: invalid channel", func(t *testing.T) {
		unregisteredChannel := "myChannel"
		invalidMsg := messageIn{
			MessageType: MessageTypeSubscribe,
			Params: paramIn{
				Channels: append(channels, channel.Channel(unregisteredChannel)),
			},
		}

		result := invalidMsg.Validate()
		expectedResult := new(validation.Result)
		expectedResult.AddFieldError(validation.FieldChannels+":"+unregisteredChannel, errorx.ErrorMsgChannelsIsEmpty)
		assert.Equal(t, expectedResult, result)
	})
}

func registerChannels() []channel.Channel {
	channelList := []string{
		"notification",
		"alert",
		"feed",
	}

	var channels []channel.Channel
	for _, ch := range channelList {
		channels = append(channels, channel.RegisterPublicChannel(ch))
	}

	return channels
}
