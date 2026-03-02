package rules

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/hebertzin/cqrs/internal/domain/entity"
)

type VelocityCounter interface {
	CountRecentByAccountID(ctx context.Context, accountID uuid.UUID, withinMinutes int) (int, error)
}

type VelocityRule struct {
	MaxPerHour int
	counter    VelocityCounter
}

func NewVelocityRule(maxPerHour int, counter VelocityCounter) *VelocityRule {
	return &VelocityRule{
		MaxPerHour: maxPerHour,
		counter:    counter,
	}
}

func (r *VelocityRule) Name() string { return "velocity" }

func (r *VelocityRule) Evaluate(ctx context.Context, tx *entity.Transaction) Result {
	count, err := r.counter.CountRecentByAccountID(ctx, tx.AccountID, 60)
	if err != nil {
		// Fail open: don't block transactions on counter errors
		return Result{}
	}

	if count >= r.MaxPerHour {
		return Result{
			Triggered: true,
			Score:     0.5,
			Reason:    fmt.Sprintf("account exceeded %d transactions per hour (current: %d)", r.MaxPerHour, count),
		}
	}

	return Result{}
}
