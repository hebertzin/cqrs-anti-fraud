package rules

import (
	"context"

	"github.com/hebertzin/cqrs/internal/domain/entity"
)

type Result struct {
	Triggered bool
	Score     float64
	Reason    string
}

type Rule interface {
	Name() string
	Evaluate(ctx context.Context, tx *entity.Transaction) Result
}

type EvaluationResult struct {
	TotalScore     float64
	Reasons        []string
	TriggeredRules []string
}

type Engine struct {
	rules []Rule
}

func NewEngine(rules ...Rule) *Engine {
	return &Engine{rules: rules}
}

func (e *Engine) Evaluate(ctx context.Context, tx *entity.Transaction) EvaluationResult {
	result := EvaluationResult{}

	for _, rule := range e.rules {
		r := rule.Evaluate(ctx, tx)
		if r.Triggered {
			result.TotalScore += r.Score
			result.Reasons = append(result.Reasons, r.Reason)
			result.TriggeredRules = append(result.TriggeredRules, rule.Name())
		}
	}

	if result.TotalScore > 1.0 {
		result.TotalScore = 1.0
	}

	return result
}
