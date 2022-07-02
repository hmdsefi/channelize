package common

type ConnectionWrapper interface {
	ID() string
	SendMessage([]byte) error
}
