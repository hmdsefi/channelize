/**
 * Copyright © 2022 Hamed Yousefi <hdyousefi@gmail.com>.
 */

package mock

type Connection struct {
	id       string
	userID   *string
	err      error
	send     chan []byte
	authFunc func() error
}

func NewConnection(id string, userID *string, authFunc func() error) *Connection {
	return &Connection{
		id:       id,
		userID:   userID,
		send:     make(chan []byte, 256),
		authFunc: authFunc,
	}
}

func (c Connection) ID() string {
	return c.id
}

func (c Connection) UserID() *string {
	return c.userID
}

func (c Connection) Authenticate() error {
	return c.authFunc()
}

func (c Connection) SendMessage(data []byte) error {
	if c.err != nil {
		return c.err
	}

	c.send <- data
	return nil
}

func (c Connection) Message() <-chan []byte {
	return c.send
}

func (c Connection) Close() {
	close(c.send)
}

func (c Connection) WithError(err error) Connection {
	conn := c
	conn.err = err
	return conn
}
