package events

import "context"

type Event interface {
	EventName() string
	EventKey() string
}

type Publisher interface {
	Publish(ctx context.Context, event Event) error
}
