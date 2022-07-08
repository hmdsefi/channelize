/**
 * Copyright © 2022 Hamed Yousefi <hdyousefi@gmail.com>.
 */

package common

// ConnectionWrapper is an interface that wraps websocket.Conn object.
type ConnectionWrapper interface {
	ID() string
	UserID() *string
	Authenticate() error
	SendMessage([]byte) error
}
