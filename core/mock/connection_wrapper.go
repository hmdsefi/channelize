/**
 * Copyright © 2022 Hamed Yousefi <hdyousefi@gmail.com>.
 */

package mock

type Connection struct {
	id   string
	send chan []byte
}

func NewConnection(id string) *Connection {
	return &Connection{
		id:   id,
		send: make(chan []byte, 256),
	}
}

func NewConnectionWithChan(id string, send chan []byte) *Connection {
	return &Connection{
		id:   id,
		send: send,
	}
}

func (c Connection) ID() string {
	return c.id
}

func (c Connection) SendMessage(data []byte) error {
	c.send <- data
	return nil
}

func (c Connection) Message() <-chan []byte {
	return c.send
}

func (c Connection) Close() {
	close(c.send)
}
