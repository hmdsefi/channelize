package main

import (
	"context"
	"time"

	uuid "github.com/satori/go.uuid"

	"github.com/hamed-yousefi/channelize/internal/channel"
)

const (
	defaultDelay = 1000 * time.Millisecond
)

type messageProducerFunc func(userID string) interface{}
type randomUserID func() string

type messageSender interface {
	SendPrivateMessage(ctx context.Context, ch channel.Channel, userID string, message interface{}) error
}

func publish(
	ctx context.Context,
	sender messageSender,
	produce messageProducerFunc,
	ch channel.Channel,
	randUserID randomUserID,
) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			userID := randUserID()
			_ = sender.SendPrivateMessage(ctx, ch, userID, produce(userID))
			time.Sleep(defaultDelay)
		}
	}
}

type Notification struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	Body      string `json:"body"`
	CreatedAt int64  `json:"created_at"`
}

func newNotification(userID string) interface{} {
	return Notification{
		ID:        uuid.NewV4().String(),
		UserID:    userID,
		Body:      "Public notification content!",
		CreatedAt: time.Now().UnixMilli(),
	}
}
