/**
 * Copyright © 2022 Hamed Yousefi <hdyousefi@gmail.com>.
 */

package core

import (
	"context"
	"encoding/json"

	"github.com/hamed-yousefi/channelize/internal/channel"
	"github.com/hamed-yousefi/channelize/internal/common"
	"github.com/hamed-yousefi/channelize/internal/common/errorx"
	"github.com/hamed-yousefi/channelize/log"
)

// store stores connections per channel.
type store interface {
	// Connections returns a list of available connections for an input channel.
	Connections(_ context.Context, ch channel.Channel) []common.ConnectionWrapper
}

// Dispatch is a mechanism to send the public and private messages to the
// available connection per channel. It uses a storage to get the connections.
type Dispatch struct {
	store  store
	logger log.Logger
}

// NewDispatch creates a new instance of Dispatch struct.
func NewDispatch(store store, logger log.Logger) *Dispatch {
	return &Dispatch{
		store:  store,
		logger: logger,
	}
}

// SendPublicMessage sends the input message to the available connections of
// input channel.
//
// This process is thread safe if the store.Connections be thread safe.
//
// SendPublicMessage might return json marshal error.
func (d *Dispatch) SendPublicMessage(ctx context.Context, ch channel.Channel, message interface{}) error {
	connections := d.store.Connections(ctx, ch)

	if len(connections) == 0 {
		return nil
	}

	msgOut := newMessageOut(ch, message)
	msgOutBytes, err := json.Marshal(msgOut)
	if err != nil {
		return errorx.NewChannelizeErrorWithErr(errorx.CodeFailedToMarshalMessage, err)
	}

	for _, connection := range connections {
		if err := connection.SendMessage(msgOutBytes); err != nil {
			d.logger.Error(
				"failed to send public message to the inbound buffer",
				common.LogFieldID, connection.ID(),
				common.LogFieldError, err.Error(),
			)
		}
	}

	return nil
}
