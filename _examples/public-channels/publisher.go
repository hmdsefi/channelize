/**
 * Copyright © 2022 Hamed Yousefi <hdyousefi@gmail.com>.
 */

package main

import (
	"context"
	"math/rand"
	"time"

	uuid "github.com/satori/go.uuid"

	"github.com/hmdsefi/channelize/internal/channel"
)

const (
	defaultDelay = 1000 * time.Millisecond
)

var (
	alertTypes      = []string{"WARNING", "ERROR", "PANIC"}
	alertPriorities = []string{"P1", "P2", "P3", "P4", "P5"}
)

type messageProducerFunc func() interface{}

type messageSender interface {
	SendPublicMessage(ctx context.Context, ch channel.Channel, message interface{}) error
}

type news struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	CreatedAt int64  `json:"created_at"`
}

func newNews() interface{} {
	return news{
		ID:        uuid.NewV4().String(),
		Title:     "News title",
		Body:      "Here you can read more about the news!",
		CreatedAt: time.Now().UnixMilli(),
	}
}

type alert struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Type        string `json:"type"`
	Priority    string `json:"priority"`
	Description string `json:"description"`
	CreatedAt   int64  `json:"created_at"`
}

func newAlert() interface{} {
	rand.Seed(time.Now().Unix())
	alertType := alertTypes[rand.Intn(len(alertTypes))]
	return alert{
		ID:          uuid.NewV4().String(),
		Title:       alertType + ": Alert title",
		Type:        alertType,
		Priority:    alertPriorities[rand.Intn(len(alertPriorities))],
		Description: "Here you can read more about the alert!",
		CreatedAt:   time.Now().UnixMilli(),
	}
}

type Notification struct {
	ID        string `json:"id"`
	Body      string `json:"body"`
	CreatedAt int64  `json:"created_at"`
}

func newNotification() interface{} {
	return Notification{
		ID:        uuid.NewV4().String(),
		Body:      "Public notification content!",
		CreatedAt: time.Now().UnixMilli(),
	}
}

func publish(ctx context.Context, sender messageSender, produce messageProducerFunc, ch channel.Channel) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			_ = sender.SendPublicMessage(ctx, ch, produce())
			time.Sleep(defaultDelay)
		}
	}
}
