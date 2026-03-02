package projector_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"github.com/hebertzin/cqrs/internal/application/eventbus"
	"github.com/hebertzin/cqrs/internal/domain/event"
	"github.com/hebertzin/cqrs/internal/infrastructure/projector"
	"github.com/hebertzin/cqrs/internal/query/model"
)

type accountBusStub struct {
	blockedHandler eventbus.Handler
}

func (b *accountBusStub) Publish(_ context.Context, _ event.Event) error { return nil }

func (b *accountBusStub) Subscribe(eventType event.Type, handler eventbus.Handler) {
	if eventType == event.TypeAccountBlocked {
		b.blockedHandler = handler
	}
}

func TestAccountProjector_AccountBlocked(t *testing.T) {
	accRepo := &mockAccReadRepo{}
	accountID := uuid.New()
	now := time.Now()

	existing := &model.AccountStatusView{
		ID:        accountID,
		Status:    "active",
		CreatedAt: now,
	}

	accRepo.On("GetByID", mock.Anything, accountID).Return(existing, nil)
	accRepo.On("Save", mock.Anything, mock.Anything).Return(nil)

	p := projector.NewAccountProjector(accRepo, zap.NewNop())
	bus := &accountBusStub{}
	p.Register(bus)

	evt := event.NewAccountBlocked(accountID, "fraud detected", "system")
	err := bus.blockedHandler(context.Background(), evt)

	assert.NoError(t, err)
	assert.Equal(t, "blocked", existing.Status)
	assert.NotNil(t, existing.BlockedAt)

	accRepo.AssertExpectations(t)
}

func TestAccountProjector_AccountNotFoundIsHandledGracefully(t *testing.T) {
	accRepo := &mockAccReadRepo{}
	accountID := uuid.New()

	accRepo.On("GetByID", mock.Anything, accountID).Return(nil, assert.AnError)

	p := projector.NewAccountProjector(accRepo, zap.NewNop())
	bus := &accountBusStub{}
	p.Register(bus)

	evt := event.NewAccountBlocked(accountID, "test", "system")
	err := bus.blockedHandler(context.Background(), evt)

	// Graceful: don't fail if account view doesn't exist yet
	assert.NoError(t, err)
}
