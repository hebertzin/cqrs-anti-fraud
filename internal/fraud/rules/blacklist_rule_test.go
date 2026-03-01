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

type mockBlacklist struct {
	mock.Mock
}

func (m *mockBlacklist) IsAccountBlacklisted(ctx context.Context, accountID uuid.UUID) (bool, error) {
	args := m.Called(ctx, accountID)
	return args.Bool(0), args.Error(1)
}

func (m *mockBlacklist) IsMerchantBlacklisted(ctx context.Context, merchantID string) (bool, error) {
	args := m.Called(ctx, merchantID)
	return args.Bool(0), args.Error(1)
}

func TestBlacklistRule_AccountBlacklisted(t *testing.T) {
	bl := &mockBlacklist{}
	accountID := uuid.New()
	bl.On("IsAccountBlacklisted", mock.Anything, accountID).Return(true, nil)

	rule := rules.NewBlacklistRule(bl)
	tx := entity.NewTransaction(accountID, 100, "BRL", "merchant", "BR")

	result := rule.Evaluate(context.Background(), tx)

	assert.True(t, result.Triggered)
	assert.Equal(t, 1.0, result.Score)
}

func TestBlacklistRule_MerchantBlacklisted(t *testing.T) {
	bl := &mockBlacklist{}
	accountID := uuid.New()
	bl.On("IsAccountBlacklisted", mock.Anything, accountID).Return(false, nil)
	bl.On("IsMerchantBlacklisted", mock.Anything, "bad-merchant").Return(true, nil)

	rule := rules.NewBlacklistRule(bl)
	tx := entity.NewTransaction(accountID, 100, "BRL", "bad-merchant", "BR")

	result := rule.Evaluate(context.Background(), tx)

	assert.True(t, result.Triggered)
	assert.Equal(t, 0.9, result.Score)
}

func TestBlacklistRule_NeitherBlacklisted(t *testing.T) {
	bl := &mockBlacklist{}
	accountID := uuid.New()
	bl.On("IsAccountBlacklisted", mock.Anything, accountID).Return(false, nil)
	bl.On("IsMerchantBlacklisted", mock.Anything, mock.Anything).Return(false, nil)

	rule := rules.NewBlacklistRule(bl)
	tx := entity.NewTransaction(accountID, 100, "BRL", "safe-merchant", "BR")

	result := rule.Evaluate(context.Background(), tx)

	assert.False(t, result.Triggered)
}

func TestBlacklistRule_AccountCheckError_FallsThrough(t *testing.T) {
	bl := &mockBlacklist{}
	accountID := uuid.New()
	bl.On("IsAccountBlacklisted", mock.Anything, accountID).Return(false, errors.New("timeout"))
	bl.On("IsMerchantBlacklisted", mock.Anything, mock.Anything).Return(false, nil)

	rule := rules.NewBlacklistRule(bl)
	tx := entity.NewTransaction(accountID, 100, "BRL", "merchant", "BR")

	result := rule.Evaluate(context.Background(), tx)

	assert.False(t, result.Triggered)
}

func TestBlacklistRule_Name(t *testing.T) {
	rule := rules.NewBlacklistRule(&mockBlacklist{})
	assert.Equal(t, "blacklist", rule.Name())
}
