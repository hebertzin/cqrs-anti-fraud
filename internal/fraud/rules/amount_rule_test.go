package rules_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/hebertzin/cqrs/internal/domain/entity"
	"github.com/hebertzin/cqrs/internal/fraud/rules"
)

func TestAmountRule_BelowThreshold(t *testing.T) {
	rule := rules.NewAmountRule(1000)
	tx := entity.NewTransaction(uuid.New(), 500, "BRL", "merchant", "BR")

	result := rule.Evaluate(context.Background(), tx)

	assert.False(t, result.Triggered)
	assert.Zero(t, result.Score)
}

func TestAmountRule_AboveThreshold(t *testing.T) {
	rule := rules.NewAmountRule(1000)
	tx := entity.NewTransaction(uuid.New(), 1500, "BRL", "merchant", "BR")

	result := rule.Evaluate(context.Background(), tx)

	assert.True(t, result.Triggered)
	assert.Equal(t, 0.4, result.Score)
	assert.NotEmpty(t, result.Reason)
}

func TestAmountRule_ExactlyThreshold(t *testing.T) {
	rule := rules.NewAmountRule(1000)
	tx := entity.NewTransaction(uuid.New(), 1000, "BRL", "merchant", "BR")

	result := rule.Evaluate(context.Background(), tx)

	assert.False(t, result.Triggered)
}

func TestAmountRule_Name(t *testing.T) {
	rule := rules.NewAmountRule(1000)
	assert.Equal(t, "amount_threshold", rule.Name())
}
