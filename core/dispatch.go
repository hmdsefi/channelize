package core

import (
	"github.com/hamed-yousefi/channelize/channel"
	"github.com/hamed-yousefi/channelize/conn"
)

type store interface {
	Connections(ch channel.Channel) []*conn.Connection
}

type Dispatch struct {
	store store
}
