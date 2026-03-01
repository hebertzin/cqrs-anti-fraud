package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/hebertzin/cqrs/internal/application/eventbus"
	"github.com/hebertzin/cqrs/internal/domain/event"
)

type EventBusMock struct {
	mock.Mock
}

func (m *EventBusMock) Publish(ctx context.Context, e event.Event) error {
	args := m.Called(ctx, e)
	return args.Error(0)
}

func (m *EventBusMock) Subscribe(eventType event.Type, handler eventbus.Handler) {
	m.Called(eventType, handler)
}
