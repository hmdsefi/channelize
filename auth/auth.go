/**
 * Copyright © 2022 Hamed Yousefi <hdyousefi@gmail.com>.
 */

package auth

// AuthenticateFunc is a function type that is responsible to authenticate
// the input token. It is an auth middleware that should be implemented by
// the client to validate user token before subscribing to a private channel
// and sending the message to a private channel.
type AuthenticateFunc func(token string) (*Token, error)

// Token represent the client websocket token details.
type Token struct {
	// Token represents client websocket token that sends it via the MessageIn.
	Token string

	// UserID represent client userID. The user that is owned the Token.
	UserID string

	// ExpiresAt represents Token expiration time. The value of ExpiresAt is
	// unix seconds.
	ExpiresAt int64
}
