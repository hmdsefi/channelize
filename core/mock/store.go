/**
 * Copyright © 2022 Hamed Yousefi <hdyousefi@gmail.com>.
 */

package mock

import (
	"context"

	"github.com/hamed-yousefi/channelize/channel"
	"github.com/hamed-yousefi/channelize/common"
)

type store interface {
	Connections(_ context.Context, ch channel.Channel) []common.ConnectionWrapper
}

type Store struct {
	connections []common.ConnectionWrapper
}

func NewStore(connections []common.ConnectionWrapper) *Store {
	return &Store{connections: connections}
}

func (s Store) Connections(_ context.Context, _ channel.Channel) []common.ConnectionWrapper {
	return s.connections
}
