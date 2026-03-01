package repository

import (
	"context"

	"github.com/google/uuid"
<<<<<<< HEAD

=======
>>>>>>> c4f71672c010347ab5a24d9bfd7962045ae3009e
	"github.com/hebertzin/cqrs/internal/domain/entity"
)

type TransactionWriteRepository interface {
	Save(ctx context.Context, transaction *entity.Transaction) error
	Update(ctx context.Context, transaction *entity.Transaction) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Transaction, error)
	FindByAccountID(ctx context.Context, accountID uuid.UUID) ([]*entity.Transaction, error)
	CountRecentByAccountID(ctx context.Context, accountID uuid.UUID, withinMinutes int) (int, error)
}
