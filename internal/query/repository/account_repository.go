package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/hebertzin/cqrs/internal/query/model"
)

type AccountReadRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*model.AccountStatusView, error)
	Save(ctx context.Context, view *model.AccountStatusView) error
	IncrementTransactionCount(ctx context.Context, accountID uuid.UUID) error
	IncrementFlaggedCount(ctx context.Context, accountID uuid.UUID) error
	IncrementDeclinedCount(ctx context.Context, accountID uuid.UUID) error
}
