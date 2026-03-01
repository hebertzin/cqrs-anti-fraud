package event_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

<<<<<<< HEAD
=======
	"github.com/hebertzin/cqrs/internal/domain/entity"
>>>>>>> c4f71672c010347ab5a24d9bfd7962045ae3009e
	"github.com/hebertzin/cqrs/internal/domain/event"
)

func TestNewBase(t *testing.T) {
<<<<<<< HEAD
	aggregateID := uuid.New()
	b := event.NewBase(event.TypeTransactionAnalyzed, aggregateID)

	assert.NotEqual(t, uuid.Nil, b.GetID())
	assert.Equal(t, event.TypeTransactionAnalyzed, b.GetType())
	assert.Equal(t, aggregateID, b.GetAggregateID())
	assert.False(t, b.GetOccurredAt().IsZero())
=======
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
>>>>>>> c4f71672c010347ab5a24d9bfd7962045ae3009e
}

func TestNewAccountBlocked(t *testing.T) {
	accountID := uuid.New()
<<<<<<< HEAD
	e := event.NewAccountBlocked(accountID, "fraud detected", "system")

	assert.Equal(t, event.TypeAccountBlocked, e.GetType())
	assert.Equal(t, accountID, e.AccountID)
	assert.Equal(t, "fraud detected", e.Reason)
	assert.Equal(t, "system", e.BlockedBy)
=======
	evt := event.NewAccountBlocked(accountID, "fraud detected", "system")

	assert.Equal(t, event.TypeAccountBlocked, evt.GetType())
	assert.Equal(t, accountID, evt.AccountID)
	assert.Equal(t, "fraud detected", evt.Reason)
	assert.Equal(t, "system", evt.BlockedBy)
>>>>>>> c4f71672c010347ab5a24d9bfd7962045ae3009e
}

func TestNewTransactionFlagged(t *testing.T) {
	txID := uuid.New()
	accID := uuid.New()
<<<<<<< HEAD
	e := event.NewTransactionFlagged(txID, accID, "suspicious pattern", "analyst")

	assert.Equal(t, event.TypeTransactionFlagged, e.GetType())
	assert.Equal(t, txID, e.TransactionID)
	assert.Equal(t, accID, e.AccountID)
	assert.Equal(t, "suspicious pattern", e.Reason)
	assert.Equal(t, "analyst", e.FlaggedBy)
=======
	evt := event.NewTransactionFlagged(txID, accID, "manual review", "analyst-1")

	assert.Equal(t, event.TypeTransactionFlagged, evt.GetType())
	assert.Equal(t, txID, evt.TransactionID)
	assert.Equal(t, accID, evt.AccountID)
	assert.Equal(t, "manual review", evt.Reason)
>>>>>>> c4f71672c010347ab5a24d9bfd7962045ae3009e
}
