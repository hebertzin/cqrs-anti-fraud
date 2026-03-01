package rules

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hebertzin/cqrs/internal/domain/entity"
)

// Blacklist provides a lookup for blocked accounts and merchants.
type Blacklist interface {
	IsAccountBlacklisted(ctx context.Context, accountID uuid.UUID) (bool, error)
	IsMerchantBlacklisted(ctx context.Context, merchantID string) (bool, error)
}

// BlacklistRule flags transactions from blacklisted accounts or merchants.
type BlacklistRule struct {
	blacklist Blacklist
}

func NewBlacklistRule(blacklist Blacklist) *BlacklistRule {
	return &BlacklistRule{blacklist: blacklist}
}

func (r *BlacklistRule) Name() string { return "blacklist" }

func (r *BlacklistRule) Evaluate(ctx context.Context, tx *entity.Transaction) Result {
	if blocked, err := r.blacklist.IsAccountBlacklisted(ctx, tx.AccountID); err == nil && blocked {
		return Result{
			Triggered: true,
			Score:     1.0,
			Reason:    fmt.Sprintf("account %s is blacklisted", tx.AccountID),
		}
	}

	if blocked, err := r.blacklist.IsMerchantBlacklisted(ctx, tx.MerchantID); err == nil && blocked {
		return Result{
			Triggered: true,
			Score:     0.9,
			Reason:    fmt.Sprintf("merchant %s is blacklisted", tx.MerchantID),
		}
	}

	return Result{}
}
