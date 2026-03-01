package inmemory

import (
	"context"
	"sync"

	"go.uber.org/zap"

	"github.com/hebertzin/cqrs/internal/application/eventbus"
	"github.com/hebertzin/cqrs/internal/domain/event"
)

type EventBus struct {
	mu       sync.RWMutex
	handlers map[event.Type][]eventbus.Handler
	logger   *zap.Logger
}

func NewEventBus(logger *zap.Logger) *EventBus {
	return &EventBus{
		handlers: make(map[event.Type][]eventbus.Handler),
		logger:   logger,
	}
}

func (b *EventBus) Publish(ctx context.Context, e event.Event) error {
	b.mu.RLock()
	handlers, ok := b.handlers[e.GetType()]
	b.mu.RUnlock()

	if !ok {
		return nil
	}

	for _, h := range handlers {
		if err := h(ctx, e); err != nil {
			b.logger.Error("event handler error",
				zap.String("event_type", string(e.GetType())),
				zap.String("event_id", e.GetID().String()),
				zap.Error(err),
			)
		}
	}
	return nil
}

func (b *EventBus) Subscribe(eventType event.Type, handler eventbus.Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[eventType] = append(b.handlers[eventType], handler)
}
