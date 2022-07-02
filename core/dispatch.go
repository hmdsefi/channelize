package core

import (
	"encoding/json"

	"github.com/hamed-yousefi/channelize/channel"
	"github.com/hamed-yousefi/channelize/common/errorx"
	"github.com/hamed-yousefi/channelize/conn"
)

type store interface {
	Connections(ch channel.Channel) []*conn.Connection
}

type Dispatch struct {
	store store
}

func NewDispatch(store store) *Dispatch {
	return &Dispatch{store: store}
}

func (d *Dispatch) SendPublicMessage(ch channel.Channel, data []byte) error {
	connections := d.store.Connections(ch)

	if len(connections) == 0 {
		return nil
	}

	msgOut := newMessageOut(ch, data)
	msgOutBytes, err := json.Marshal(msgOut)
	if err != nil {
		return errorx.NewChannelizeErrorWithErr(errorx.CodeFailedToMarshalMessage, err)
	}

	for _, connection := range connections {
		if err := connection.SendMessage(msgOutBytes); err != nil {
			// TODO log the error
		}
	}

	return nil
}
