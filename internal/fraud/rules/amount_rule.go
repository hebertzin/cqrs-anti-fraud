package rules

import (
	"context"
	"fmt"

	"github.com/hebertzin/cqrs/internal/domain/entity"
)

// AmountRule flags transactions that exceed a configurable threshold.
type AmountRule struct {
	Threshold float64
}

func NewAmountRule(threshold float64) *AmountRule {
	return &AmountRule{Threshold: threshold}
}

func (r *AmountRule) Name() string { return "amount_threshold" }

func (r *AmountRule) Evaluate(_ context.Context, tx *entity.Transaction) Result {
	if tx.Amount > r.Threshold {
		return Result{
			Triggered: true,
			Score:     0.4,
			Reason:    fmt.Sprintf("transaction amount %.2f exceeds threshold %.2f", tx.Amount, r.Threshold),
		}
	}
	return Result{}
}
