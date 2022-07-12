/**
 * Copyright © 2022 Hamed Yousefi <hdyousefi@gmail.com>.
 */

package mock

import (
	"context"

	"github.com/hmdsefi/channelize/internal/channel"
	"github.com/hmdsefi/channelize/internal/common"
)

type Store struct {
	connections []common.ConnectionWrapper
	send        chan string
}

func NewStore(connections []common.ConnectionWrapper) *Store {
	return &Store{
		connections: connections,
		send:        make(chan string, 1),
	}
}

func (s Store) UnsubscribeUserID(_ context.Context, _ string, _ string, ch channel.Channel) {
	s.send <- ch.String()
}

func (s Store) Connections(_ context.Context, _ channel.Channel) []common.ConnectionWrapper {
	return s.connections
}

func (s Store) ConnectionByUserID(_ context.Context, _ channel.Channel, _ string) common.ConnectionWrapper {
	if len(s.connections) > 0 {
		return s.connections[0]
	}

	return nil
}

func (s Store) Receive() string {
	return <-s.send
}
