package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/hebertzin/cqrs/internal/query/model"
)

type TransactionReadRepositoryMock struct {
	mock.Mock
}

func (m *TransactionReadRepositoryMock) GetByID(ctx context.Context, id uuid.UUID) (*model.TransactionRiskView, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.TransactionRiskView), args.Error(1)
}

func (m *TransactionReadRepositoryMock) GetByAccountID(ctx context.Context, accountID uuid.UUID) ([]*model.TransactionRiskView, error) {
	args := m.Called(ctx, accountID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.TransactionRiskView), args.Error(1)
}

func (m *TransactionReadRepositoryMock) Save(ctx context.Context, view *model.TransactionRiskView) error {
	args := m.Called(ctx, view)
	return args.Error(0)
}

func (m *TransactionReadRepositoryMock) GetFraudAlerts(ctx context.Context, page, limit int) (*model.FraudAlertListResponse, error) {
	args := m.Called(ctx, page, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.FraudAlertListResponse), args.Error(1)
}
