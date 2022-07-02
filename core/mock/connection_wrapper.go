package mock

type Connection struct {
	id string
}

func NewConnection(id string) *Connection {
	return &Connection{id: id}
}

func (c Connection) ID() string {
	return c.id
}

func (c Connection) SendMessage(_ []byte) error {
	return nil
}
