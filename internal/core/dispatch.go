/**
 * Copyright © 2022 Hamed Yousefi <hdyousefi@gmail.com>.
 */

package core

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/hmdsefi/channelize/internal/channel"
	"github.com/hmdsefi/channelize/internal/common"
	"github.com/hmdsefi/channelize/internal/common/errorx"
	"github.com/hmdsefi/channelize/log"
)

// store stores connections per channel.
type store interface {
	// UnsubscribeByUserID removes the input channel subscription from the storage. Also,
	// removes userID and connID mapping to prevent token validation for each input message.
	UnsubscribeUserID(_ context.Context, connID string, userID string, channels channel.Channel)

	// Connections returns a list of available connections for an input channel.
	Connections(ctx context.Context, ch channel.Channel) []common.ConnectionWrapper

	// ConnectionByUserID returns a connection that mapped with input userID and channel.
	ConnectionByUserID(ctx context.Context, ch channel.Channel, userID string) common.ConnectionWrapper
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

	for _, conn := range connections {
		if err := conn.SendMessage(msgOutBytes); err != nil {
			d.logger.Error(
				"failed to send public message to the inbound buffer",
				common.LogFieldID, conn.ID(),
				common.LogFieldError, err.Error(),
			)
		}
	}

	return nil
}

// SendPrivateMessage sends the input message to the input channel if the client
// already authenticated with the input userID. Otherwise, skips and returns.
//
// On each call it will authenticate the connection to ensure token is not expired.
//
// SendPrivateMessage might return token expiration or json marshal errors.
func (d *Dispatch) SendPrivateMessage(ctx context.Context, ch channel.Channel, userID string, message interface{}) error {
	conn := d.store.ConnectionByUserID(ctx, ch, userID)
	if conn == nil {
		return nil
	}

	// validate auth token before sending the message.
	err := conn.Authenticate()
	var authErr *errorx.ChannelizeError
	switch {
	case err == nil:
	case errors.As(err, &authErr):
		if authErr.Code == errorx.CodeAuthTokenIsMissing ||
			authErr.Code == errorx.CodeAuthFuncIsMissing ||
			authErr.Code == errorx.CodeAuthTokenIsExpired {
			d.store.UnsubscribeUserID(ctx, conn.ID(), userID, ch)
		}
		// TODO write error to the connection
		return err
	default:
		// TODO write error to the connection
		return err
	}

	msgOut := newMessageOut(ch, message)
	msgOutBytes, err := json.Marshal(msgOut)
	if err != nil {
		return errorx.NewChannelizeErrorWithErr(errorx.CodeFailedToMarshalMessage, err)
	}

	if err = conn.SendMessage(msgOutBytes); err != nil {
		d.logger.Error(
			"failed to send private message to the inbound buffer",
			common.LogFieldID, conn.ID(),
			common.LogFieldError, err.Error(),
		)

		return err
	}

	return nil
}
