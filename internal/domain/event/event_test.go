package event_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/hebertzin/cqrs/internal/domain/event"
)

func TestNewBase(t *testing.T) {
	aggregateID := uuid.New()
	b := event.NewBase(event.TypeTransactionAnalyzed, aggregateID)

	assert.NotEqual(t, uuid.Nil, b.GetID())
	assert.Equal(t, event.TypeTransactionAnalyzed, b.GetType())
	assert.Equal(t, aggregateID, b.GetAggregateID())
	assert.False(t, b.GetOccurredAt().IsZero())
}

func TestNewAccountBlocked(t *testing.T) {
	accountID := uuid.New()
	e := event.NewAccountBlocked(accountID, "fraud detected", "system")

	assert.Equal(t, event.TypeAccountBlocked, e.GetType())
	assert.Equal(t, accountID, e.AccountID)
	assert.Equal(t, "fraud detected", e.Reason)
	assert.Equal(t, "system", e.BlockedBy)
}

func TestNewTransactionFlagged(t *testing.T) {
	txID := uuid.New()
	accID := uuid.New()
	e := event.NewTransactionFlagged(txID, accID, "suspicious pattern", "analyst")

	assert.Equal(t, event.TypeTransactionFlagged, e.GetType())
	assert.Equal(t, txID, e.TransactionID)
	assert.Equal(t, accID, e.AccountID)
	assert.Equal(t, "suspicious pattern", e.Reason)
	assert.Equal(t, "analyst", e.FlaggedBy)
}
