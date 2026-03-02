package handler_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	queryhandler "github.com/hebertzin/cqrs/internal/query/handler"
	"github.com/hebertzin/cqrs/internal/query/model"
)

type mockAccountReadRepo struct {
	mock.Mock
}

func (m *mockAccountReadRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.AccountStatusView, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.AccountStatusView), args.Error(1)
}

func (m *mockAccountReadRepo) Save(ctx context.Context, view *model.AccountStatusView) error {
	return m.Called(ctx, view).Error(0)
}

func (m *mockAccountReadRepo) IncrementTransactionCount(ctx context.Context, accountID uuid.UUID) error {
	return m.Called(ctx, accountID).Error(0)
}

func (m *mockAccountReadRepo) IncrementFlaggedCount(ctx context.Context, accountID uuid.UUID) error {
	return m.Called(ctx, accountID).Error(0)
}

func (m *mockAccountReadRepo) IncrementDeclinedCount(ctx context.Context, accountID uuid.UUID) error {
	return m.Called(ctx, accountID).Error(0)
}

func TestGetAccountStatusHandler_Found(t *testing.T) {
	readRepo := &mockAccountReadRepo{}
	accountID := uuid.New()

	expected := &model.AccountStatusView{
		ID:     accountID,
		Status: "active",
	}
	readRepo.On("GetByID", mock.Anything, accountID).Return(expected, nil)

	handler := queryhandler.NewGetAccountStatusHandler(readRepo)
	result, err := handler.Handle(context.Background(), queryhandler.GetAccountStatusQuery{AccountID: accountID})

	assert.NoError(t, err)
	view := result.(*model.AccountStatusView)
	assert.Equal(t, accountID, view.ID)
	assert.Equal(t, "active", view.Status)
}

func TestGetAccountStatusHandler_NotFound(t *testing.T) {
	readRepo := &mockAccountReadRepo{}
	readRepo.On("GetByID", mock.Anything, mock.Anything).Return(nil, errors.New("not found"))

	handler := queryhandler.NewGetAccountStatusHandler(readRepo)
	_, err := handler.Handle(context.Background(), queryhandler.GetAccountStatusQuery{AccountID: uuid.New()})

	assert.Error(t, err)
}

func TestGetAccountStatusHandler_InvalidQuery(t *testing.T) {
	handler := queryhandler.NewGetAccountStatusHandler(&mockAccountReadRepo{})
	_, err := handler.Handle(context.Background(), "not-a-query")
	assert.Error(t, err)
}
