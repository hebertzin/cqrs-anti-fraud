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
	"github.com/hebertzin/cqrs/internal/domain/entity"
	"github.com/hebertzin/cqrs/internal/domain/event"
	"github.com/hebertzin/cqrs/internal/infrastructure/projector"
	"github.com/hebertzin/cqrs/internal/query/model"
)

type mockTxReadRepo struct{ mock.Mock }

func (m *mockTxReadRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.TransactionRiskView, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.TransactionRiskView), args.Error(1)
}

func (m *mockTxReadRepo) GetByAccountID(ctx context.Context, accountID uuid.UUID) ([]*model.TransactionRiskView, error) {
	args := m.Called(ctx, accountID)
	return nil, args.Error(1)
}

func (m *mockTxReadRepo) Save(ctx context.Context, view *model.TransactionRiskView) error {
	return m.Called(ctx, view).Error(0)
}

func (m *mockTxReadRepo) GetFraudAlerts(ctx context.Context, page, limit int) (*model.FraudAlertListResponse, error) {
	return nil, nil
}

type mockAccReadRepo struct{ mock.Mock }

func (m *mockAccReadRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.AccountStatusView, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.AccountStatusView), args.Error(1)
}

func (m *mockAccReadRepo) Save(ctx context.Context, view *model.AccountStatusView) error {
	return m.Called(ctx, view).Error(0)
}

func (m *mockAccReadRepo) IncrementTransactionCount(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

func (m *mockAccReadRepo) IncrementFlaggedCount(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

func (m *mockAccReadRepo) IncrementDeclinedCount(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

// --- Tests ---

func makeTransactionAnalyzedEvent(riskScore float64, status entity.TransactionStatus) event.TransactionAnalyzed {
	tx := entity.NewTransaction(uuid.New(), 100, "BRL", "merchant", "BR")
	tx.SetRiskScore(riskScore)
	tx.Status = status
	return event.NewTransactionAnalyzed(tx, []string{"test reason"})
}

func TestTransactionProjector_LowRiskEvent(t *testing.T) {
	txRepo := &mockTxReadRepo{}
	accRepo := &mockAccReadRepo{}

	txRepo.On("Save", mock.Anything, mock.Anything).Return(nil)
	accRepo.On("IncrementTransactionCount", mock.Anything, mock.Anything).Return(nil)

	p := projector.NewTransactionProjector(txRepo, accRepo, zap.NewNop())
	bus := &eventBusStub{}
	p.Register(bus)

	evt := makeTransactionAnalyzedEvent(0.2, entity.TransactionStatusApproved)
	err := bus.publish(context.Background(), evt)

	assert.NoError(t, err)
	txRepo.AssertExpectations(t)
	accRepo.AssertExpectations(t)
}

func TestTransactionProjector_HighRiskEvent(t *testing.T) {
	txRepo := &mockTxReadRepo{}
	accRepo := &mockAccReadRepo{}

	txRepo.On("Save", mock.Anything, mock.Anything).Return(nil)
	accRepo.On("IncrementTransactionCount", mock.Anything, mock.Anything).Return(nil)
	accRepo.On("IncrementFlaggedCount", mock.Anything, mock.Anything).Return(nil)
	accRepo.On("IncrementDeclinedCount", mock.Anything, mock.Anything).Return(nil)

	p := projector.NewTransactionProjector(txRepo, accRepo, zap.NewNop())
	bus := &eventBusStub{}
	p.Register(bus)

	evt := makeTransactionAnalyzedEvent(0.9, entity.TransactionStatusDeclined)
	err := bus.publish(context.Background(), evt)

	assert.NoError(t, err)
	txRepo.AssertExpectations(t)
	accRepo.AssertExpectations(t)
}

func TestTransactionProjector_FlaggedEvent(t *testing.T) {
	txID := uuid.New()
	txRepo := &mockTxReadRepo{}
	accRepo := &mockAccReadRepo{}

	existing := &model.TransactionRiskView{
		ID:        txID,
		AccountID: uuid.New(),
		Status:    "approved",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	txRepo.On("GetByID", mock.Anything, txID).Return(existing, nil)
	txRepo.On("Save", mock.Anything, mock.Anything).Return(nil)

	p := projector.NewTransactionProjector(txRepo, accRepo, zap.NewNop())
	bus := &eventBusStub{}
	p.Register(bus)

	evt := event.NewTransactionFlagged(txID, uuid.New(), "manual review", "analyst")
	err := bus.publishFlagged(context.Background(), evt)

	assert.NoError(t, err)
	txRepo.AssertExpectations(t)
}

// eventBusStub routes published events to the right handler without a real bus.
type eventBusStub struct {
	analyzedHandler eventbus.Handler
	flaggedHandler  eventbus.Handler
}

func (b *eventBusStub) Publish(_ context.Context, _ event.Event) error { return nil }

func (b *eventBusStub) Subscribe(eventType event.Type, handler eventbus.Handler) {
	switch eventType {
	case event.TypeTransactionAnalyzed:
		b.analyzedHandler = handler
	case event.TypeTransactionFlagged:
		b.flaggedHandler = handler
	}
}

func (b *eventBusStub) publish(ctx context.Context, e event.Event) error {
	if b.analyzedHandler != nil {
		return b.analyzedHandler(ctx, e)
	}
	return nil
}

func (b *eventBusStub) publishFlagged(ctx context.Context, e event.Event) error {
	if b.flaggedHandler != nil {
		return b.flaggedHandler(ctx, e)
	}
	return nil
}
