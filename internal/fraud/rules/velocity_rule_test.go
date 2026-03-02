package rules_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/hebertzin/cqrs/internal/domain/entity"
	"github.com/hebertzin/cqrs/internal/fraud/rules"
)

type mockVelocityCounter struct {
	mock.Mock
}

func (m *mockVelocityCounter) CountRecentByAccountID(ctx context.Context, accountID uuid.UUID, withinMinutes int) (int, error) {
	args := m.Called(ctx, accountID, withinMinutes)
	return args.Int(0), args.Error(1)
}

func TestVelocityRule_BelowLimit(t *testing.T) {
	counter := &mockVelocityCounter{}
	counter.On("CountRecentByAccountID", mock.Anything, mock.Anything, 60).Return(3, nil)

	rule := rules.NewVelocityRule(10, counter)
	tx := entity.NewTransaction(uuid.New(), 100, "BRL", "merchant", "BR")

	result := rule.Evaluate(context.Background(), tx)

	assert.False(t, result.Triggered)
	counter.AssertExpectations(t)
}

func TestVelocityRule_AtLimit(t *testing.T) {
	counter := &mockVelocityCounter{}
	counter.On("CountRecentByAccountID", mock.Anything, mock.Anything, 60).Return(10, nil)

	rule := rules.NewVelocityRule(10, counter)
	tx := entity.NewTransaction(uuid.New(), 100, "BRL", "merchant", "BR")

	result := rule.Evaluate(context.Background(), tx)

	assert.True(t, result.Triggered)
	assert.Equal(t, 0.5, result.Score)
	assert.Contains(t, result.Reason, "10")
}

func TestVelocityRule_CounterError_FailsOpen(t *testing.T) {
	counter := &mockVelocityCounter{}
	counter.On("CountRecentByAccountID", mock.Anything, mock.Anything, 60).Return(0, errors.New("redis down"))

	rule := rules.NewVelocityRule(5, counter)
	tx := entity.NewTransaction(uuid.New(), 100, "BRL", "merchant", "BR")

	result := rule.Evaluate(context.Background(), tx)

	// Should fail open: don't block on counter errors
	assert.False(t, result.Triggered)
}

func TestVelocityRule_Name(t *testing.T) {
	rule := rules.NewVelocityRule(10, &mockVelocityCounter{})
	assert.Equal(t, "velocity", rule.Name())
}
