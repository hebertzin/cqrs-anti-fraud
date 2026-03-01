package repository

import (
	"context"

	"github.com/google/uuid"
<<<<<<< HEAD

=======
>>>>>>> c4f71672c010347ab5a24d9bfd7962045ae3009e
	"github.com/hebertzin/cqrs/internal/domain/entity"
)

type AccountWriteRepository interface {
	Save(ctx context.Context, account *entity.Account) error
	Update(ctx context.Context, account *entity.Account) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Account, error)
}
