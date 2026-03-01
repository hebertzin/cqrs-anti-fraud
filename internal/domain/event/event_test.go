package event_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/hebertzin/cqrs/internal/domain/entity"
	"github.com/hebertzin/cqrs/internal/domain/event"
)

func TestNewBase(t *testing.T) {
	aggID := uuid.New()
	e := event.NewBase(event.TypeTransactionAnalyzed, aggID)

	assert.NotEqual(t, uuid.Nil, e.GetID())
	assert.Equal(t, event.TypeTransactionAnalyzed, e.GetType())
	assert.Equal(t, aggID, e.GetAggregateID())
	assert.False(t, e.GetOccurredAt().IsZero())
}

func TestNewTransactionAnalyzed(t *testing.T) {
	tx := entity.NewTransaction(uuid.New(), 500, "BRL", "merchant", "BR")
	tx.SetRiskScore(0.7)
	tx.Flag()
	reasons := []string{"amount exceeded", "suspicious location"}

	evt := event.NewTransactionAnalyzed(tx, reasons)

	assert.Equal(t, event.TypeTransactionAnalyzed, evt.GetType())
	assert.Equal(t, tx.ID, evt.TransactionID)
	assert.Equal(t, tx.AccountID, evt.AccountID)
	assert.Equal(t, tx.Amount, evt.Amount)
	assert.Equal(t, tx.RiskScore, evt.RiskScore)
	assert.Equal(t, tx.Status, evt.Status)
	assert.Equal(t, reasons, evt.FraudReasons)
}

func TestNewAccountBlocked(t *testing.T) {
	accountID := uuid.New()
	evt := event.NewAccountBlocked(accountID, "fraud detected", "system")

	assert.Equal(t, event.TypeAccountBlocked, evt.GetType())
	assert.Equal(t, accountID, evt.AccountID)
	assert.Equal(t, "fraud detected", evt.Reason)
	assert.Equal(t, "system", evt.BlockedBy)
}

func TestNewTransactionFlagged(t *testing.T) {
	txID := uuid.New()
	accID := uuid.New()
	evt := event.NewTransactionFlagged(txID, accID, "manual review", "analyst-1")

	assert.Equal(t, event.TypeTransactionFlagged, evt.GetType())
	assert.Equal(t, txID, evt.TransactionID)
	assert.Equal(t, accID, evt.AccountID)
	assert.Equal(t, "manual review", evt.Reason)
}
