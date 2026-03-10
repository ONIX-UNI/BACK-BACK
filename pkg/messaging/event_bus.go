package messaging

import "context"

type DomainEvent interface {
	EventName() string
}

type EventBus interface {
	Publish(ctx context.Context, event DomainEvent) error
}
