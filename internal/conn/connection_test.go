/**
 * Copyright © 2022 Hamed Yousefi <hdyousefi@gmail.com>.
 */

package conn

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hamed-yousefi/channelize/auth"
	"github.com/hamed-yousefi/channelize/internal/common/errorx"
	"github.com/hamed-yousefi/channelize/internal/common/log"
	"github.com/hamed-yousefi/channelize/internal/common/utils"
)

const (
	protocolHTTP = "http"
	protocolWS   = "ws"
	wsPath       = "/ws"

	testAuthToken = "test-auth-token" // nolint
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

func newMockHelper(receive chan<- string) *MockMessageProcessor {
	return &MockMessageProcessor{
		receive: receive,
	}
}

func (m MockMessageProcessor) Remove(_ context.Context, _ string, _ *string) {
}

func (m MockMessageProcessor) ParseMessage(_ context.Context, _ *Connection, message []byte) {
	m.receive <- string(message)
}

func (m MockMessageProcessor) close() {
	close(m.receive)
}

func testAuthenticateFunc(token string) (*auth.Token, error) {
	return &auth.Token{}, nil
}

type connStore struct {
	mu          sync.RWMutex
	connections []*Connection
}

func (c *connStore) add(conn *Connection) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.connections = append(c.connections, conn)
}

func (c *connStore) get(index int) *Connection {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.connections[index]
}

func (c *connStore) len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.connections)
}

type Handler struct {
	http.Handler
	ctx              context.Context
	cancel           context.CancelFunc
	mockMsgProcessor *MockMessageProcessor
	connStore        *connStore
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
		connStore:        new(connStore),
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
		s.connStore.add(NewConnection(
			s.ctx, conn,
			s.mockMsgProcessor,
			testAuthenticateFunc,
			log.NewDefaultLogger(),
		))
	}
}

// TestNewConnection tests Connection initialization and read/write methods.
func TestNewConnection(t *testing.T) {
	receiver := make(chan string)
	mockMsgProcessor := newMockHelper(receiver)
	defer mockMsgProcessor.close()

	handler := newHandler(t, mockMsgProcessor)
	server := httptest.NewServer(handler)
	defer server.Close()

	wsURL := protocolWS + strings.TrimPrefix(server.URL, protocolHTTP) + wsPath

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
		assert.Equal(t, expectedClientMsg, actualClientMsg)
	})

	t.Run("Test write to client", func(t *testing.T) {
		expectedServerMsg := "test write message"

		require.Equal(t, 1, handler.connStore.len())

		err = handler.connStore.get(0).SendMessage([]byte(expectedServerMsg))
		require.Nil(t, err)

		msgType, msg, err := ws.ReadMessage()
		if err != nil {
			t.Fatalf("failed to read a message websocket connection, %v", err)
		}

		require.Equal(t, websocket.TextMessage, msgType)

		actualServerMsg := string(msg)
		assert.Equal(t, expectedServerMsg, actualServerMsg)
	})

	_ = handler.Close()
}

func TestConnection_UserID(t *testing.T) {
	t.Run("nil token", func(t *testing.T) {
		conn := Connection{}
		assert.Nil(t, conn.UserID())
	})

	t.Run("empty userID", func(t *testing.T) {
		conn := Connection{
			token: &auth.Token{},
		}
		assert.Nil(t, conn.UserID())
	})

	t.Run("valid userID", func(t *testing.T) {
		expectedUserID := "test-user-id"
		conn := Connection{
			token: &auth.Token{
				UserID: expectedUserID,
			},
		}

		actualUserID := conn.UserID()
		require.NotNil(t, actualUserID)
		assert.Equal(t, expectedUserID, *actualUserID)
	})
}

func TestConnection_AuthenticateAndStore(t *testing.T) {
	t.Run("nil authFunc", func(t *testing.T) {
		t.Parallel()
		conn := Connection{}
		err := conn.AuthenticateAndStore(testAuthToken)
		require.NotNil(t, err)
		var chanErr *errorx.ChannelizeError
		require.True(t, errors.As(err, &chanErr))
		assert.Equal(t, errorx.CodeAuthFuncIsMissing, chanErr.Code)
		assert.Equal(t, errorx.ErrorMsgAuthFuncIsMissing, chanErr.Error())
	})

	t.Run("validation error", func(t *testing.T) {
		t.Parallel()
		errMsg := "invalid token"
		authFunc := func(_ string) (*auth.Token, error) {
			return nil, errors.New(errMsg)
		}
		conn := Connection{authFunc: authFunc}
		err := conn.AuthenticateAndStore(testAuthToken)
		require.NotNil(t, err)
		assert.Equal(t, errMsg, err.Error())
	})

	t.Run("expired token", func(t *testing.T) {
		t.Parallel()
		authFunc := func(_ string) (*auth.Token, error) {
			return &auth.Token{
				ExpiresAt: utils.Now().Add(-1 * time.Minute).Unix(),
			}, nil
		}
		conn := Connection{authFunc: authFunc}
		err := conn.AuthenticateAndStore(testAuthToken)
		require.NotNil(t, err)
		var chanErr *errorx.ChannelizeError
		require.True(t, errors.As(err, &chanErr))
		assert.Equal(t, errorx.CodeAuthTokenIsExpired, chanErr.Code)
		assert.Equal(t, errorx.ErrorMsgAuthTokenIsExpired, chanErr.Error())
	})

	t.Run("valid token", func(t *testing.T) {
		t.Parallel()
		token := auth.Token{
			ExpiresAt: utils.Now().Add(time.Minute).Unix(),
		}

		authFunc := func(_ string) (*auth.Token, error) {
			out := token
			return &out, nil
		}
		conn := Connection{authFunc: authFunc}
		err := conn.AuthenticateAndStore(testAuthToken)
		require.Nil(t, err)
		assert.Equal(t, token, *conn.token)
	})
}

func TestConnection_Authenticate(t *testing.T) {
	t.Run("nil authFunc", func(t *testing.T) {
		t.Parallel()
		conn := Connection{}
		err := conn.Authenticate()
		require.NotNil(t, err)
		var chanErr *errorx.ChannelizeError
		require.True(t, errors.As(err, &chanErr))
		assert.Equal(t, errorx.CodeAuthTokenIsMissing, chanErr.Code)
		assert.Equal(t, errorx.ErrorMsgConnectionAuthTokenIsMissing, chanErr.Error())
	})

	t.Run("expired token", func(t *testing.T) {
		t.Parallel()
		token := auth.Token{
			ExpiresAt: utils.Now().Add(-1 * time.Minute).Unix(),
		}

		authFunc := func(_ string) (*auth.Token, error) {
			out := token
			return &out, nil
		}

		conn := Connection{authFunc: authFunc, token: &token}
		err := conn.Authenticate()
		require.NotNil(t, err)
		var chanErr *errorx.ChannelizeError
		require.True(t, errors.As(err, &chanErr))
		assert.Equal(t, errorx.CodeAuthTokenIsExpired, chanErr.Code)
		assert.Equal(t, errorx.ErrorMsgAuthTokenIsExpired, chanErr.Error())
	})

	t.Run("valid token", func(t *testing.T) {
		t.Parallel()
		token := auth.Token{
			ExpiresAt: utils.Now().Add(time.Minute).Unix(),
		}

		authFunc := func(_ string) (*auth.Token, error) {
			out := token
			return &out, nil
		}

		conn := Connection{authFunc: authFunc, token: &token}
		err := conn.Authenticate()
		assert.Nil(t, err)
	})

	t.Run("extended token", func(t *testing.T) {
		t.Parallel()
		token := auth.Token{
			ExpiresAt: utils.Now().Add(-1 * time.Minute).Unix(),
		}

		extendedToken := auth.Token{
			ExpiresAt: utils.Now().Add(time.Minute).Unix(),
		}

		authFunc := func(_ string) (*auth.Token, error) {
			out := extendedToken
			return &out, nil
		}

		conn := Connection{authFunc: authFunc, token: &token}
		err := conn.Authenticate()
		assert.Nil(t, err)
		assert.Equal(t, extendedToken, *conn.token)
	})
}

// TODO test concurrent connections
// TODO test context cancellation propagation for multiple connections
// TODO test server side closing connection
// TODO test client side closing connection
// TODO test ping/pong failures
