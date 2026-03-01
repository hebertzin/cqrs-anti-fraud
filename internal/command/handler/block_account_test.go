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
	"github.com/hebertzin/cqrs/tests/mocks"
	"github.com/hebertzin/cqrs/pkg/logger"
)

func TestBlockAccountHandler_Success(t *testing.T) {
	accRepo := &mocks.AccountWriteRepositoryMock{}
	eventBus := &mocks.EventBusMock{}
	log := logger.NewNop()

	accountID := uuid.New()
	account := entity.NewAccount(uuid.New())
	account.ID = accountID

	accRepo.On("FindByID", mock.Anything, accountID).Return(account, nil)
	accRepo.On("Update", mock.Anything, mock.Anything).Return(nil)
	eventBus.On("Publish", mock.Anything, mock.Anything).Return(nil)

	handler := cmdhandler.NewBlockAccountHandler(accRepo, eventBus, log)
	result, err := handler.Handle(context.Background(), cmdmodel.BlockAccount{
		AccountID: accountID,
		Reason:    "suspicious activity",
		BlockedBy: "fraud-system",
	})

	assert.NoError(t, err)
	assert.Nil(t, result)
	assert.True(t, account.IsBlocked())

	accRepo.AssertExpectations(t)
	eventBus.AssertExpectations(t)
}

func TestBlockAccountHandler_AlreadyBlocked(t *testing.T) {
	accRepo := &mocks.AccountWriteRepositoryMock{}
	eventBus := &mocks.EventBusMock{}
	log := logger.NewNop()

	accountID := uuid.New()
	account := entity.NewAccount(uuid.New())
	account.ID = accountID
	account.Block()

	accRepo.On("FindByID", mock.Anything, accountID).Return(account, nil)

	handler := cmdhandler.NewBlockAccountHandler(accRepo, eventBus, log)
	_, err := handler.Handle(context.Background(), cmdmodel.BlockAccount{
		AccountID: accountID,
		Reason:    "test",
		BlockedBy: "operator",
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already blocked")
}

func TestBlockAccountHandler_AccountNotFound(t *testing.T) {
	accRepo := &mocks.AccountWriteRepositoryMock{}
	log := logger.NewNop()

	accountID := uuid.New()
	accRepo.On("FindByID", mock.Anything, accountID).Return(nil, errors.New("not found"))

	handler := cmdhandler.NewBlockAccountHandler(accRepo, &mocks.EventBusMock{}, log)
	_, err := handler.Handle(context.Background(), cmdmodel.BlockAccount{AccountID: accountID})

	assert.Error(t, err)
}

func TestBlockAccountHandler_InvalidCommand(t *testing.T) {
	handler := cmdhandler.NewBlockAccountHandler(nil, nil, logger.NewNop())
	_, err := handler.Handle(context.Background(), "invalid")
	assert.Error(t, err)
}
