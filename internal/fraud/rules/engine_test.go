package rules_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/hebertzin/cqrs/internal/domain/entity"
	"github.com/hebertzin/cqrs/internal/fraud/rules"
)

func TestEngine_NoRulesTriggered(t *testing.T) {
	tx := entity.NewTransaction(uuid.New(), 100, "BRL", "merchant", "BR")
	engine := rules.NewEngine(rules.NewAmountRule(10000), rules.NewLocationRule(nil))

	result := engine.Evaluate(context.Background(), tx)

	assert.Zero(t, result.TotalScore)
	assert.Empty(t, result.Reasons)
	assert.Empty(t, result.TriggeredRules)
}

func TestEngine_CapsScoreAt1(t *testing.T) {
	tx := entity.NewTransaction(uuid.New(), 99999, "BRL", "merchant", "XX")
	engine := rules.NewEngine(
		rules.NewAmountRule(100),
		rules.NewLocationRule(nil),
	)

	result := engine.Evaluate(context.Background(), tx)

	assert.LessOrEqual(t, result.TotalScore, 1.0)
	assert.Len(t, result.TriggeredRules, 2)
}

func TestEngine_MultipleRules(t *testing.T) {
	tx := entity.NewTransaction(uuid.New(), 50000, "BRL", "merchant", "XX")
	engine := rules.NewEngine(
		rules.NewAmountRule(1000),
		rules.NewLocationRule(nil),
	)

	result := engine.Evaluate(context.Background(), tx)

	assert.Greater(t, result.TotalScore, 0.0)
	assert.Len(t, result.Reasons, 2)
}
