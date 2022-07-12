package main

import (
	"encoding/hex"
	"math/rand"
	"time"

	uuid "github.com/satori/go.uuid"

	"github.com/hmdsefi/channelize/internal/common/utils"
)

const (
	defaultTokenTTL = 30 * time.Minute
)

type Token struct {
	UserID    string `json:"id"`
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
}

func NewToken(token string) *Token {
	return &Token{
		UserID:    uuid.NewV4().String(),
		Token:     token,
		ExpiresAt: utils.Now().Add(defaultTokenTTL).Unix(),
	}
}

func (t *Token) ExtendTTL() {
	t.ExpiresAt = utils.Now().Add(defaultTokenTTL).Unix()
}

type authService struct {
	tokens  map[string]*Token
	userIDs []string
}

func newAuth() *authService {
	return &authService{
		tokens: make(map[string]*Token),
	}
}

func (a *authService) CreateToken() string {
	tokenBytes := make([]byte, 30)
	rand.Read(tokenBytes) // nolint
	tokenStr := hex.EncodeToString(tokenBytes)
	a.tokens[tokenStr] = NewToken(tokenStr)
	a.userIDs = append(a.userIDs, a.tokens[tokenStr].UserID)
	return tokenStr
}

func (a *authService) Token(token string) *Token {
	return a.tokens[token]
}
