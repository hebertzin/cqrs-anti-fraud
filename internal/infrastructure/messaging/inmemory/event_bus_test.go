package inmemory_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/hebertzin/cqrs/internal/domain/event"
	"github.com/hebertzin/cqrs/internal/infrastructure/messaging/inmemory"
)

func TestEventBus_PublishWithNoSubscribers(t *testing.T) {
	bus := inmemory.NewEventBus(zap.NewNop())
	e := event.NewBase(event.TypeTransactionAnalyzed, [16]byte{})
	err := bus.Publish(context.Background(), e)
	assert.NoError(t, err)
}

func TestEventBus_PublishCallsSubscriber(t *testing.T) {
	bus := inmemory.NewEventBus(zap.NewNop())
	called := false

	bus.Subscribe(event.TypeTransactionAnalyzed, func(_ context.Context, _ event.Event) error {
		called = true
		return nil
	})

	e := event.NewBase(event.TypeTransactionAnalyzed, [16]byte{})
	err := bus.Publish(context.Background(), e)

	assert.NoError(t, err)
	assert.True(t, called)
}

func TestEventBus_PublishCallsMultipleSubscribers(t *testing.T) {
	bus := inmemory.NewEventBus(zap.NewNop())
	count := 0

	inc := func(_ context.Context, _ event.Event) error {
		count++
		return nil
	}

	bus.Subscribe(event.TypeTransactionAnalyzed, inc)
	bus.Subscribe(event.TypeTransactionAnalyzed, inc)

	e := event.NewBase(event.TypeTransactionAnalyzed, [16]byte{})
	bus.Publish(context.Background(), e) //nolint:errcheck

	assert.Equal(t, 2, count)
}

func TestEventBus_SubscriberErrorIsLoggedNotReturned(t *testing.T) {
	bus := inmemory.NewEventBus(zap.NewNop())

	bus.Subscribe(event.TypeAccountBlocked, func(_ context.Context, _ event.Event) error {
		return assert.AnError
	})

	e := event.NewBase(event.TypeAccountBlocked, [16]byte{})
	err := bus.Publish(context.Background(), e)

	// Publish itself should not fail even if handler errors
	assert.NoError(t, err)
}

func TestEventBus_DifferentEventTypesIsolated(t *testing.T) {
	bus := inmemory.NewEventBus(zap.NewNop())
	txCalled := false
	accCalled := false

	bus.Subscribe(event.TypeTransactionAnalyzed, func(_ context.Context, _ event.Event) error {
		txCalled = true
		return nil
	})
	bus.Subscribe(event.TypeAccountBlocked, func(_ context.Context, _ event.Event) error {
		accCalled = true
		return nil
	})

	txEvent := event.NewBase(event.TypeTransactionAnalyzed, [16]byte{})
	bus.Publish(context.Background(), txEvent) //nolint:errcheck

	assert.True(t, txCalled)
	assert.False(t, accCalled)
}
