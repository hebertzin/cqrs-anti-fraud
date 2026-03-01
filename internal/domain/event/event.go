package event

import (
	"time"

	"github.com/google/uuid"
)

type Type string

const (
	TypeTransactionAnalyzed Type = "transaction.analyzed"
	TypeTransactionFlagged  Type = "transaction.flagged"
	TypeAccountBlocked      Type = "account.blocked"
	TypeAccountActivated    Type = "account.activated"
)

type Event interface {
	GetID() uuid.UUID
	GetType() Type
	GetAggregateID() uuid.UUID
	GetOccurredAt() time.Time
}

type Base struct {
	ID          uuid.UUID `json:"id"`
<<<<<<< HEAD
	EventType   Type      `json:"event_type"`
=======
	EventType   Type      `json:"type"`
>>>>>>> c4f71672c010347ab5a24d9bfd7962045ae3009e
	AggregateID uuid.UUID `json:"aggregate_id"`
	OccurredAt  time.Time `json:"occurred_at"`
}

func NewBase(eventType Type, aggregateID uuid.UUID) Base {
	return Base{
		ID:          uuid.New(),
		EventType:   eventType,
		AggregateID: aggregateID,
		OccurredAt:  time.Now().UTC(),
	}
}

<<<<<<< HEAD
func (b Base) GetID() uuid.UUID          { return b.ID }
func (b Base) GetType() Type             { return b.EventType }
func (b Base) GetAggregateID() uuid.UUID { return b.AggregateID }
func (b Base) GetOccurredAt() time.Time  { return b.OccurredAt }
=======
func (e Base) GetID() uuid.UUID        { return e.ID }
func (e Base) GetType() Type           { return e.EventType }
func (e Base) GetAggregateID() uuid.UUID { return e.AggregateID }
func (e Base) GetOccurredAt() time.Time  { return e.OccurredAt }
>>>>>>> c4f71672c010347ab5a24d9bfd7962045ae3009e
