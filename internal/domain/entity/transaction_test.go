package entity_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/hebertzin/cqrs/internal/domain/entity"
)

func TestNewTransaction(t *testing.T) {
	accountID := uuid.New()
	tx := entity.NewTransaction(accountID, 100.50, "BRL", "merchant-1", "BR")

	assert.NotEqual(t, uuid.Nil, tx.ID)
	assert.Equal(t, accountID, tx.AccountID)
	assert.Equal(t, 100.50, tx.Amount)
	assert.Equal(t, "BRL", tx.Currency)
	assert.Equal(t, "merchant-1", tx.MerchantID)
	assert.Equal(t, "BR", tx.Location)
	assert.Equal(t, entity.TransactionStatusPending, tx.Status)
	assert.Zero(t, tx.RiskScore)
	assert.False(t, tx.CreatedAt.IsZero())
}

func TestTransaction_Approve(t *testing.T) {
	tx := entity.NewTransaction(uuid.New(), 50, "BRL", "m1", "BR")
	tx.Approve()
	assert.Equal(t, entity.TransactionStatusApproved, tx.Status)
}

func TestTransaction_Decline(t *testing.T) {
	tx := entity.NewTransaction(uuid.New(), 50, "BRL", "m1", "BR")
	tx.Decline()
	assert.Equal(t, entity.TransactionStatusDeclined, tx.Status)
}

func TestTransaction_Flag(t *testing.T) {
	tx := entity.NewTransaction(uuid.New(), 50, "BRL", "m1", "BR")
	tx.Flag()
	assert.Equal(t, entity.TransactionStatusFlagged, tx.Status)
}

func TestTransaction_SetRiskScore(t *testing.T) {
	tx := entity.NewTransaction(uuid.New(), 50, "BRL", "m1", "BR")
	tx.SetRiskScore(0.75)
	assert.Equal(t, 0.75, tx.RiskScore)
}

func TestTransaction_IsHighRisk(t *testing.T) {
	tx := entity.NewTransaction(uuid.New(), 50, "BRL", "m1", "BR")

	tx.SetRiskScore(0.9)
	assert.True(t, tx.IsHighRisk())
	assert.False(t, tx.IsMediumRisk())

	tx.SetRiskScore(0.6)
	assert.False(t, tx.IsHighRisk())
	assert.True(t, tx.IsMediumRisk())

	tx.SetRiskScore(0.3)
	assert.False(t, tx.IsHighRisk())
	assert.False(t, tx.IsMediumRisk())
}
