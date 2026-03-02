package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/hebertzin/cqrs/internal/domain/entity"
)

type TransactionWriteRepositoryMock struct {
	mock.Mock
}

func (m *TransactionWriteRepositoryMock) Save(ctx context.Context, tx *entity.Transaction) error {
	args := m.Called(ctx, tx)
	return args.Error(0)
}

func (m *TransactionWriteRepositoryMock) Update(ctx context.Context, tx *entity.Transaction) error {
	args := m.Called(ctx, tx)
	return args.Error(0)
}

func (m *TransactionWriteRepositoryMock) FindByID(ctx context.Context, id uuid.UUID) (*entity.Transaction, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Transaction), args.Error(1)
}

func (m *TransactionWriteRepositoryMock) FindByAccountID(ctx context.Context, accountID uuid.UUID) ([]*entity.Transaction, error) {
	args := m.Called(ctx, accountID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Transaction), args.Error(1)
}

func (m *TransactionWriteRepositoryMock) CountRecentByAccountID(ctx context.Context, accountID uuid.UUID, withinMinutes int) (int, error) {
	args := m.Called(ctx, accountID, withinMinutes)
	return args.Int(0), args.Error(1)
}
