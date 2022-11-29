package core

type Sender interface {
	SendMessage(message string) bool
}
