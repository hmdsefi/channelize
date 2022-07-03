/**
 * Copyright © 2022 Hamed Yousefi <hdyousefi@gmail.com>.
 */

package conn

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hamed-yousefi/channelize/common/validation"
)

const (
	protocolHTTP = "http"
	protocolWS   = "ws"
	wsPath       = "/ws"
)

var (
	wsUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

type MockMessageProcessor struct {
	receive chan<- string
}

func newMockMessageProcessor(receive chan<- string) *MockMessageProcessor {
	return &MockMessageProcessor{
		receive: receive,
	}
}

func (m MockMessageProcessor) Validate(_ []byte) validation.Result {
	return validation.Result{}
}

func (m MockMessageProcessor) ProcessMessage(_ context.Context, _ *Connection, message []byte) {
	m.receive <- string(message)
}

func (m MockMessageProcessor) close() {
	close(m.receive)
}

type Handler struct {
	http.Handler
	ctx              context.Context
	cancel           context.CancelFunc
	mockMsgProcessor *MockMessageProcessor
	connections      []*Connection
}

func (s *Handler) Close() error {
	s.cancel()
	return nil
}

func newHandler(t *testing.T, mockMsgProcessor *MockMessageProcessor) *Handler {
	ctx, cancel := context.WithCancel(context.Background())

	svr := &Handler{
		ctx:              ctx,
		cancel:           cancel,
		mockMsgProcessor: mockMsgProcessor,
	}

	router := http.NewServeMux()
	router.Handle(wsPath, svr.makeWebsocketHandler(t))

	svr.Handler = router

	return svr
}

func (s *Handler) makeWebsocketHandler(t *testing.T) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := wsUpgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatal(err)
		}
		s.connections = append(s.connections, NewConnection(s.ctx, conn, s.mockMsgProcessor))
	}
}

// TestNewConnection tests Connection initialization and read/write methods.
func TestNewConnection(t *testing.T) {
	receiver := make(chan string)
	mockMsgProcessor := newMockMessageProcessor(receiver)
	defer mockMsgProcessor.close()

	handler := newHandler(t, mockMsgProcessor)
	server := httptest.NewServer(handler)
	defer server.Close()

	wsURL := protocolWS + strings.TrimPrefix(server.URL, protocolHTTP) + wsPath
	t.Log("ws URL:", wsURL)

	ws, resp, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("failed to open websocket connection, %s %v", wsURL, err)
	}
	defer func() {
		_ = resp.Body.Close()
		_ = ws.Close()
	}()

	t.Run("Test read from client", func(t *testing.T) {
		expectedClientMsg := "test read message"
		if err = ws.WriteMessage(websocket.TextMessage, []byte(expectedClientMsg)); err != nil {
			t.Fatalf("failed to send a message over websocket connection, %v", err)
		}

		actualClientMsg := <-receiver
		t.Log("message that receive by the server:", actualClientMsg)
		assert.Equal(t, expectedClientMsg, actualClientMsg)
	})

	t.Run("Test write to client", func(t *testing.T) {
		expectedServerMsg := "test write message"

		require.Equal(t, 1, len(handler.connections))

		err = handler.connections[0].SendMessage([]byte(expectedServerMsg))
		require.Nil(t, err)

		msgType, msg, err := ws.ReadMessage()
		if err != nil {
			t.Fatalf("failed to read a message websocket connection, %v", err)
		}

		require.Equal(t, websocket.TextMessage, msgType)

		actualServerMsg := string(msg)
		t.Log("message that client read:", actualServerMsg)
		assert.Equal(t, expectedServerMsg, actualServerMsg)
	})

	_ = handler.Close()
}

// TODO test concurrent connections
// TODO test context cancellation propagation for multiple connections
// TODO test server side closing connection
// TODO test client side closing connection
// TODO test ping/pong failures
