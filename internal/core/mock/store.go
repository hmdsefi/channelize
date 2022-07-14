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
	userConnections map[string]common.ConnectionWrapper
	connections     []common.ConnectionWrapper
	send            chan string
}

func NewStore(userConnections map[string]common.ConnectionWrapper) *Store {
	var connections []common.ConnectionWrapper
	for _, conn := range userConnections {
		connections = append(connections, conn)
	}
	return &Store{
		connections:     connections,
		userConnections: userConnections,
		send:            make(chan string, 1),
	}
}

func (s Store) UnsubscribeUserID(_ context.Context, _ string, _ string, ch channel.Channel) {
	s.send <- ch.String()
}

func (s Store) Connections(_ context.Context, _ channel.Channel) []common.ConnectionWrapper {
	return s.connections
}

func (s Store) ConnectionByUserID(_ context.Context, _ channel.Channel, userID string) common.ConnectionWrapper {
	return s.userConnections[userID]
}

func (s Store) Receive() string {
	return <-s.send
}
