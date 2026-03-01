package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/hebertzin/cqrs/internal/domain/entity"
)

// TransactionWriteRepository defines write operations for transactions (command side).
type TransactionWriteRepository interface {
	Save(ctx context.Context, transaction *entity.Transaction) error
	Update(ctx context.Context, transaction *entity.Transaction) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Transaction, error)
	FindByAccountID(ctx context.Context, accountID uuid.UUID) ([]*entity.Transaction, error)
	CountRecentByAccountID(ctx context.Context, accountID uuid.UUID, withinMinutes int) (int, error)
}
