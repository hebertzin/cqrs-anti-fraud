package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/hebertzin/cqrs/internal/query/model"
)

type TransactionReadRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*model.TransactionRiskView, error)
	GetByAccountID(ctx context.Context, accountID uuid.UUID) ([]*model.TransactionRiskView, error)
	Save(ctx context.Context, view *model.TransactionRiskView) error
	GetFraudAlerts(ctx context.Context, page, limit int) (*model.FraudAlertListResponse, error)
}
