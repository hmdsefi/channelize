package channelize

import (
	"context"

	"github.com/hamed-yousefi/channelize/channel"
	"github.com/hamed-yousefi/channelize/conn"
	"github.com/hamed-yousefi/channelize/core"
)

type store interface {
	Subscribe(ctx context.Context, conn *conn.Connection, channels ...channel.Channel)
	Unsubscribe(ctx context.Context, connID string, channels ...channel.Channel)
	Remove(ctx context.Context, connID string)
}

type helper struct {
	store store
}

func newHelper(store store) *helper {
	return &helper{store: store}
}

func (h *helper) ProcessMessage(ctx context.Context, connection *conn.Connection, data []byte) {
	msg, err := core.UnmarshalMessageIn(data)
	if err != nil {
		// TODO write error to the websocket connection 'error' channel
		return
	}

	if res := msg.Validate(); !res.IsValid() {
		// TODO write error to the websocket connection 'error' channel
		return
	}

	switch msg.MessageType {
	case core.MessageTypeSubscribe:
		h.store.Subscribe(ctx, connection, msg.Params.Channels...)
	case core.MessageTypeUnsubscribe:
		h.store.Unsubscribe(ctx, connection.ID(), msg.Params.Channels...)
	}
}

func (h *helper) Remove(ctx context.Context, connID string) {
	h.store.Remove(ctx, connID)
}
