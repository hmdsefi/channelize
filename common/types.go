/**
 * Copyright © 2022 Hamed Yousefi <hdyousefi@gmail.com>.
 */

package common

type ConnectionWrapper interface {
	ID() string
	SendMessage([]byte) error
}
