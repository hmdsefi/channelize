/**
 * Copyright © 2022 Hamed Yousefi <hdyousefi@gmail.com>.
 */

package errorx

import "fmt"

const (
	CodeConnectionClosed     = 1000
	CodeOutboundBufferIsFull = 1001

	CodeFailedToUnmarshalMessage = 1500
	CodeFailedToMarshalMessage   = 1501

	CodeAuthFuncIsMissing  = 2000
	CodeAuthTokenIsMissing = 2001
	CodeAuthTokenIsExpired = 2002
)

const (
	ErrorMsgConnectionClosed             = "websocket connection is closed"
	ErrorMsgOutboundBufferIsFull         = "connection outbound buffer is full"
	ErrorMsgUnmarshalInboundMessage      = "failed to unmarshal inbound message"
	ErrorMsgMarshalOutboundMessage       = "failed to marshal outbound message"
	ErrorMsgUnsupportedMessageType       = "message type is not supported"
	ErrorMsgChannelsIsEmpty              = "channels list is empty, minimum size is 1"
	ErrorMsgUnsupportedChannel           = "channel is not supported"
	ErrorMsgInvalidChannelType           = "channel should be either private or public"
	ErrorMsgAuthTokenIsMissing           = "auth token is missing for the private channel" // nolint
	ErrorMsgFailedToCloseConnection      = "failed to close connection"
	ErrorMsgFailedToSetReadDeadline      = "failed to set read deadline"
	ErrorMsgAuthFuncIsMissing            = "authentication function to validate private auth token"
	ErrorMsgConnectionAuthTokenIsMissing = "connection auth token is nil"
	ErrorMsgAuthTokenIsExpired           = "auth token is expired" // nolint
)

var (
	code2ErrMsg = map[int]string{
		CodeConnectionClosed:         ErrorMsgConnectionClosed,
		CodeOutboundBufferIsFull:     ErrorMsgOutboundBufferIsFull,
		CodeFailedToUnmarshalMessage: ErrorMsgUnmarshalInboundMessage,
		CodeFailedToMarshalMessage:   ErrorMsgMarshalOutboundMessage,
		CodeAuthFuncIsMissing:        ErrorMsgAuthFuncIsMissing,
		CodeAuthTokenIsMissing:       ErrorMsgConnectionAuthTokenIsMissing,
		CodeAuthTokenIsExpired:       ErrorMsgAuthTokenIsExpired,
	}
)

// ChannelizeError represents a custom error object that holds
// error details.
type ChannelizeError struct {
	Code    int
	message string
	err     error
}

// NewChannelizeError creates a new custom error by using an error Code.
func NewChannelizeError(code int) *ChannelizeError {
	return &ChannelizeError{
		Code:    code,
		message: code2ErrMsg[code],
	}
}

// NewChannelizeErrorWithErr creates a new custom error by wrapping an
// existing error.
func NewChannelizeErrorWithErr(code int, err error) *ChannelizeError {
	return &ChannelizeError{
		Code:    code,
		message: code2ErrMsg[code],
		err:     err,
	}
}

func (e ChannelizeError) Error() string {
	if e.err != nil {
		return fmt.Sprintf("%s: %s", e.message, e.err)
	}

	return e.message
}
