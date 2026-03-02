package handler_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	cmdhandler "github.com/hebertzin/cqrs/internal/command/handler"
	cmdmodel "github.com/hebertzin/cqrs/internal/command/model"
	"github.com/hebertzin/cqrs/internal/domain/entity"
	"github.com/hebertzin/cqrs/pkg/logger"
	"github.com/hebertzin/cqrs/tests/mocks"
)

func TestFlagTransactionHandler_Success(t *testing.T) {
	txRepo := &mocks.TransactionWriteRepositoryMock{}
	eventBus := &mocks.EventBusMock{}
	log := logger.NewNop()

	txID := uuid.New()
	tx := entity.NewTransaction(uuid.New(), 100, "BRL", "m1", "BR")
	tx.ID = txID
	tx.Approve()

	txRepo.On("FindByID", mock.Anything, txID).Return(tx, nil)
	txRepo.On("Update", mock.Anything, mock.Anything).Return(nil)
	eventBus.On("Publish", mock.Anything, mock.Anything).Return(nil)

	handler := cmdhandler.NewFlagTransactionHandler(txRepo, eventBus, log)
	result, err := handler.Handle(context.Background(), cmdmodel.FlagTransaction{
		TransactionID: txID,
		Reason:        "manual review",
		FlaggedBy:     "analyst",
	})

	assert.NoError(t, err)
	assert.Nil(t, result)
	assert.Equal(t, entity.TransactionStatusFlagged, tx.Status)

	txRepo.AssertExpectations(t)
	eventBus.AssertExpectations(t)
}

func TestFlagTransactionHandler_NotFound(t *testing.T) {
	txRepo := &mocks.TransactionWriteRepositoryMock{}
	log := logger.NewNop()

	txID := uuid.New()
	txRepo.On("FindByID", mock.Anything, txID).Return(nil, errors.New("not found"))

	handler := cmdhandler.NewFlagTransactionHandler(txRepo, &mocks.EventBusMock{}, log)
	_, err := handler.Handle(context.Background(), cmdmodel.FlagTransaction{TransactionID: txID})

	assert.Error(t, err)
}

func TestFlagTransactionHandler_InvalidCommand(t *testing.T) {
	handler := cmdhandler.NewFlagTransactionHandler(nil, nil, logger.NewNop())
	_, err := handler.Handle(context.Background(), "invalid")
	assert.Error(t, err)
}
