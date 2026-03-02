package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/hebertzin/cqrs/internal/domain/entity"
)

type AccountWriteRepositoryMock struct {
	mock.Mock
}

func (m *AccountWriteRepositoryMock) Save(ctx context.Context, account *entity.Account) error {
	args := m.Called(ctx, account)
	return args.Error(0)
}

func (m *AccountWriteRepositoryMock) Update(ctx context.Context, account *entity.Account) error {
	args := m.Called(ctx, account)
	return args.Error(0)
}

func (m *AccountWriteRepositoryMock) FindByID(ctx context.Context, id uuid.UUID) (*entity.Account, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Account), args.Error(1)
}
