package entity_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/hebertzin/cqrs/internal/domain/entity"
)

func TestNewAccount(t *testing.T) {
	userID := uuid.New()
<<<<<<< HEAD
	acc := entity.NewAccount(userID)

	assert.NotEqual(t, uuid.Nil, acc.ID)
	assert.Equal(t, userID, acc.UserID)
	assert.Equal(t, entity.AccountStatusActive, acc.Status)
	assert.Nil(t, acc.BlockedAt)
	assert.False(t, acc.CreatedAt.IsZero())
}

func TestAccount_Block(t *testing.T) {
	acc := entity.NewAccount(uuid.New())
	acc.Block()

	assert.Equal(t, entity.AccountStatusBlocked, acc.Status)
	assert.NotNil(t, acc.BlockedAt)
	assert.True(t, acc.IsBlocked())
}

func TestAccount_Flag(t *testing.T) {
	acc := entity.NewAccount(uuid.New())
	acc.Flag()

	assert.Equal(t, entity.AccountStatusFlagged, acc.Status)
	assert.False(t, acc.IsBlocked())
}

func TestAccount_Activate(t *testing.T) {
	acc := entity.NewAccount(uuid.New())
	acc.Block()
	acc.Activate()

	assert.Equal(t, entity.AccountStatusActive, acc.Status)
	assert.Nil(t, acc.BlockedAt)
	assert.False(t, acc.IsBlocked())
}

func TestAccount_UpdateRiskLevel(t *testing.T) {
	acc := entity.NewAccount(uuid.New())
	acc.UpdateRiskLevel(0.85)

	assert.Equal(t, 0.85, acc.RiskLevel)
=======
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
>>>>>>> c4f71672c010347ab5a24d9bfd7962045ae3009e
}
