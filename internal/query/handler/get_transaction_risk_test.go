package handler_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	queryhandler "github.com/hebertzin/cqrs/internal/query/handler"
	"github.com/hebertzin/cqrs/internal/query/model"
	"github.com/hebertzin/cqrs/tests/mocks"
)

func TestGetTransactionRiskHandler_Found(t *testing.T) {
	readRepo := &mocks.TransactionReadRepositoryMock{}
	txID := uuid.New()

	expected := &model.TransactionRiskView{
		ID:        txID,
		AccountID: uuid.New(),
		Amount:    500,
		Currency:  "BRL",
		RiskScore: 0.2,
		Status:    "approved",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	readRepo.On("GetByID", mock.Anything, txID).Return(expected, nil)

	handler := queryhandler.NewGetTransactionRiskHandler(readRepo)
	result, err := handler.Handle(context.Background(), queryhandler.GetTransactionRiskQuery{TransactionID: txID})

	assert.NoError(t, err)
	view := result.(*model.TransactionRiskView)
	assert.Equal(t, txID, view.ID)
	assert.Equal(t, model.RiskLevelLow, view.RiskLevel)

	readRepo.AssertExpectations(t)
}

func TestGetTransactionRiskHandler_NotFound(t *testing.T) {
	readRepo := &mocks.TransactionReadRepositoryMock{}
	txID := uuid.New()

	readRepo.On("GetByID", mock.Anything, txID).Return(nil, errors.New("not found"))

	handler := queryhandler.NewGetTransactionRiskHandler(readRepo)
	result, err := handler.Handle(context.Background(), queryhandler.GetTransactionRiskQuery{TransactionID: txID})

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetTransactionRiskHandler_InvalidQuery(t *testing.T) {
	handler := queryhandler.NewGetTransactionRiskHandler(&mocks.TransactionReadRepositoryMock{})
	_, err := handler.Handle(context.Background(), "not-a-query")
	assert.Error(t, err)
}
