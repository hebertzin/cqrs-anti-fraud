package rules

import (
	"context"
	"fmt"

	"github.com/hebertzin/cqrs/internal/domain/entity"
)

var DefaultSuspiciousLocations = map[string]bool{
	"XX": true,
	"ZZ": true,
}

type LocationRule struct {
	suspiciousLocations map[string]bool
}

func NewLocationRule(locations map[string]bool) *LocationRule {
	if locations == nil {
		locations = DefaultSuspiciousLocations
	}
	return &LocationRule{suspiciousLocations: locations}
}

func (r *LocationRule) Name() string { return "suspicious_location" }

func (r *LocationRule) Evaluate(_ context.Context, tx *entity.Transaction) Result {
	if r.suspiciousLocations[tx.Location] {
		return Result{
			Triggered: true,
			Score:     0.5,
			Reason:    fmt.Sprintf("transaction from suspicious location: %s", tx.Location),
		}
	}
	return Result{}
}
