package eventbus

import (
	"context"

	"github.com/hebertzin/cqrs/internal/domain/event"
)

// Handler is a function that handles a domain event.
type Handler func(ctx context.Context, e event.Event) error

// EventBus is the contract for publishing and subscribing to domain events.
type EventBus interface {
	Publish(ctx context.Context, e event.Event) error
	Subscribe(eventType event.Type, handler Handler)
}
