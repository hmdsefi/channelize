package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/websocket"

	"github.com/hmdsefi/channelize"
	"github.com/hmdsefi/channelize/auth"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	authSvc := newAuth()
	chlz := channelize.NewChannelize(channelize.WithAuthFunc(makeAuthFunc(authSvc)))
	rand.Seed(time.Now().Unix())

	notificationChannel := channelize.RegisterPrivateChannel("notifications")
	go publish(ctx, chlz, newNotification, notificationChannel, func() string {
		if len(authSvc.userIDs) == 0 {
			return ""
		}

		n := rand.Intn(len(authSvc.userIDs))
		return authSvc.userIDs[n]
	})

	go initiateAndServe(ctx, authSvc, chlz)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	cancel()
}

func initiateAndServe(ctx context.Context, authSvc *authService, chlz *channelize.Channelize) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	http.HandleFunc("/ws", chlz.MakeHTTPHandler(ctx, upgrader))
	http.HandleFunc("/ws/token", makeCreateWSTokenHandler(authSvc))

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func makeAuthFunc(svc *authService) auth.AuthenticateFunc {
	return func(token string) (*auth.Token, error) {
		authToken := svc.Token(token)
		if authToken == nil {
			return nil, errors.New("invalid token")
		}

		return &auth.Token{
			UserID:    authToken.UserID,
			Token:     authToken.Token,
			ExpiresAt: authToken.ExpiresAt,
		}, nil
	}
}

type createTokenResponse struct {
	Token string `json:"token"`
}

func makeCreateWSTokenHandler(authSvc *authService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := createTokenResponse{
			Token: authSvc.CreateToken(),
		}

		data, _ := json.Marshal(resp)
		w.Write(data)   // nolint
		fmt.Fprintln(w) // nolint
	}
}
