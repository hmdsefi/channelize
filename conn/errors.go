/**
 * Copyright © 2022 Hamed Yousefi <hdyousefi@gmail.com>.
 */

package conn

const (
	CodeConnectionClosed     = 1000
	CodeOutboundBufferIsFull = 1001
)

const (
	ErrorMsgConnectionClosed     = "websocket connection is closed"
	ErrorMsgOutboundBufferIsFull = "connection outbound buffer is full"
)

var (
	code2ErrMsg = map[int]string{
		CodeConnectionClosed:     ErrorMsgConnectionClosed,
		CodeOutboundBufferIsFull: ErrorMsgOutboundBufferIsFull,
	}
)

// ChannelizeError represents a custom error object that holds
// error details.
type ChannelizeError struct {
	code    int
	message string
}

func newConnectionError(code int) *ChannelizeError {
	return &ChannelizeError{
		code:    code,
		message: code2ErrMsg[code],
	}
}

func (e ChannelizeError) Error() string {
	return e.message
}
