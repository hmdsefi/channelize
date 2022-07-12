/**
 * Copyright © 2022 Hamed Yousefi <hdyousefi@gmail.com>.
 */

package channelize

import (
	"context"

	"github.com/hmdsefi/channelize/internal/channel"
	"github.com/hmdsefi/channelize/internal/common"
	"github.com/hmdsefi/channelize/internal/conn"
	"github.com/hmdsefi/channelize/internal/core"
)

// store is an interface this provides the ability of storing mapping
// between connections and channels. It can register a connection into
// the storage or remove it from the storage.
type store interface {
	// Subscribe creates a mapping between the connection and input channels.
	Subscribe(ctx context.Context, conn common.ConnectionWrapper, channels ...channel.Channel)

	// Unsubscribe removes the existing mapping between the input connection
	// and channels.
	Unsubscribe(ctx context.Context, connID string, channels ...channel.Channel)

	// Remove removes all the subscriptions for the input connection.
	Remove(ctx context.Context, connID string, userID *string)
}

// helper provides functionalities to the connection to register and unregister
// itself into the storage.
type helper struct {
	store store
}

func newHelper(store store) *helper {
	return &helper{
		store: store,
	}
}

// ParseMessage deserializes the inbound messages and calls the storage
// methods based on the input action.
//
// It also validates the inbound messages and publishes the errors to
// the error channel.
//
// If client message contains auth token, it validates the token. If token
// is valid, it adds the token object to the client's connection. Otherwise,
// publishes the validation error to the error channel.
func (h *helper) ParseMessage(ctx context.Context, connection *conn.Connection, data []byte) {
	msg, err := core.UnmarshalMessageIn(data)
	if err != nil {
		// TODO write error to the websocket connection 'error' channel
		return
	}

	if res := msg.Validate(); !res.IsValid() {
		// TODO write error to the websocket connection 'error' channel
		return
	}

	// validate token and store it in connection if it exists in the message.
	if msg.Params.HasToken() {
		if err := connection.AuthenticateAndStore(*msg.Params.Token); err != nil {
			// TODO write error to the websocket connection 'error' channel
			return
		}
	}

	switch msg.MessageType {
	case core.MessageTypeSubscribe:
		h.store.Subscribe(ctx, connection, msg.Params.Channels...)
	case core.MessageTypeUnsubscribe:
		h.store.Unsubscribe(ctx, connection.ID(), msg.Params.Channels...)
	}
}

// Remove removes a connection from the storage.
func (h *helper) Remove(ctx context.Context, connID string, userID *string) {
	h.store.Remove(ctx, connID, userID)
}
