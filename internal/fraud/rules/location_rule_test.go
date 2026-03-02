package rules_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/hebertzin/cqrs/internal/domain/entity"
	"github.com/hebertzin/cqrs/internal/fraud/rules"
)

func TestLocationRule_SafeLocation(t *testing.T) {
	rule := rules.NewLocationRule(nil)
	tx := entity.NewTransaction(uuid.New(), 100, "BRL", "merchant", "BR")

	result := rule.Evaluate(context.Background(), tx)

	assert.False(t, result.Triggered)
}

func TestLocationRule_SuspiciousLocation(t *testing.T) {
	rule := rules.NewLocationRule(nil)
	tx := entity.NewTransaction(uuid.New(), 100, "BRL", "merchant", "XX")

	result := rule.Evaluate(context.Background(), tx)

	assert.True(t, result.Triggered)
	assert.Equal(t, 0.5, result.Score)
	assert.Contains(t, result.Reason, "XX")
}

func TestLocationRule_CustomLocations(t *testing.T) {
	customLocations := map[string]bool{"RO": true, "MD": true}
	rule := rules.NewLocationRule(customLocations)

	safe := entity.NewTransaction(uuid.New(), 100, "BRL", "merchant", "BR")
	assert.False(t, rule.Evaluate(context.Background(), safe).Triggered)

	risky := entity.NewTransaction(uuid.New(), 100, "BRL", "merchant", "RO")
	assert.True(t, rule.Evaluate(context.Background(), risky).Triggered)
}

func TestLocationRule_Name(t *testing.T) {
	rule := rules.NewLocationRule(nil)
	assert.Equal(t, "suspicious_location", rule.Name())
}
