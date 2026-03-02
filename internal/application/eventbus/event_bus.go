package eventbus

import (
	"context"

	"github.com/hebertzin/cqrs/internal/domain/event"
)

type Handler func(ctx context.Context, e event.Event) error

type EventBus interface {
	Publish(ctx context.Context, e event.Event) error
	Subscribe(eventType event.Type, handler Handler)
}
