package model_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hebertzin/cqrs/internal/query/model"
)

func TestRiskLevelFromScore(t *testing.T) {
	cases := []struct {
		score    float64
		expected model.RiskLevel
	}{
		{0.0, model.RiskLevelLow},
		{0.3, model.RiskLevelLow},
		{0.49, model.RiskLevelLow},
		{0.5, model.RiskLevelMedium},
		{0.7, model.RiskLevelMedium},
		{0.79, model.RiskLevelMedium},
		{0.8, model.RiskLevelHigh},
		{0.9, model.RiskLevelHigh},
		{1.0, model.RiskLevelHigh},
	}

	for _, c := range cases {
		assert.Equal(t, c.expected, model.RiskLevelFromScore(c.score), "score=%.2f", c.score)
	}
}
