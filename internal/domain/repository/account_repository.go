package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/hebertzin/cqrs/internal/domain/entity"
)

type AccountWriteRepository interface {
	Save(ctx context.Context, account *entity.Account) error
	Update(ctx context.Context, account *entity.Account) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Account, error)
}
