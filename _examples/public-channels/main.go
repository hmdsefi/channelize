/**
 * Copyright © 2022 Hamed Yousefi <hdyousefi@gmail.com>.
 */

package main

import (
	"context"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/websocket"
	"github.com/hmdsefi/channelize"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	chlz := channelize.NewChannelize()

	newsChannel := channelize.RegisterPublicChannel("news")
	go publish(ctx, chlz, newNews, newsChannel)

	alertsChannel := channelize.RegisterPublicChannel("alerts")
	go publish(ctx, chlz, newAlert, alertsChannel)

	notificationChannel := channelize.RegisterPublicChannel("notifications")
	go publish(ctx, chlz, newNotification, notificationChannel)

	go initiateAndServe(ctx, chlz)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	cancel()
}

func initiateAndServe(ctx context.Context, chlz *channelize.Channelize) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	http.HandleFunc("/ws", chlz.MakeHTTPHandler(ctx, upgrader))
	http.HandleFunc("/metrics", promhttp.Handler().ServeHTTP)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
