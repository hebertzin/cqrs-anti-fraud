package entity_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/hebertzin/cqrs/internal/domain/entity"
)

func TestNewAccount(t *testing.T) {
	userID := uuid.New()
	account := entity.NewAccount(userID)

	assert.NotEqual(t, uuid.Nil, account.ID)
	assert.Equal(t, userID, account.UserID)
	assert.Equal(t, entity.AccountStatusActive, account.Status)
	assert.Nil(t, account.BlockedAt)
	assert.False(t, account.CreatedAt.IsZero())
}

func TestAccount_Block(t *testing.T) {
	account := entity.NewAccount(uuid.New())
	account.Block()

	assert.Equal(t, entity.AccountStatusBlocked, account.Status)
	assert.NotNil(t, account.BlockedAt)
	assert.True(t, account.IsBlocked())
}

func TestAccount_Flag(t *testing.T) {
	account := entity.NewAccount(uuid.New())
	account.Flag()
	assert.Equal(t, entity.AccountStatusFlagged, account.Status)
	assert.False(t, account.IsBlocked())
}

func TestAccount_Activate(t *testing.T) {
	account := entity.NewAccount(uuid.New())
	account.Block()
	assert.True(t, account.IsBlocked())

	account.Activate()
	assert.Equal(t, entity.AccountStatusActive, account.Status)
	assert.Nil(t, account.BlockedAt)
	assert.False(t, account.IsBlocked())
}

func TestAccount_UpdateRiskLevel(t *testing.T) {
	account := entity.NewAccount(uuid.New())
	account.UpdateRiskLevel(0.85)
	assert.Equal(t, 0.85, account.RiskLevel)
}
