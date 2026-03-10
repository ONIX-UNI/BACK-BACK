package messaging

type Event interface {
	EventName() string
	EventKey() string
}
