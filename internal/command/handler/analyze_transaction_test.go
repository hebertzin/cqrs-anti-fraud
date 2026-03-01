package handler_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	cmdhandler "github.com/hebertzin/cqrs/internal/command/handler"
	cmdmodel "github.com/hebertzin/cqrs/internal/command/model"
	"github.com/hebertzin/cqrs/internal/fraud/rules"
	"github.com/hebertzin/cqrs/tests/mocks"
	"github.com/hebertzin/cqrs/pkg/logger"
)

func TestAnalyzeTransactionHandler_LowRisk(t *testing.T) {
	txRepo := &mocks.TransactionWriteRepositoryMock{}
	accRepo := &mocks.AccountWriteRepositoryMock{}
	eventBus := &mocks.EventBusMock{}
	log := logger.NewNop()

	txRepo.On("Save", mock.Anything, mock.Anything).Return(nil)
	txRepo.On("CountRecentByAccountID", mock.Anything, mock.Anything, mock.Anything).Return(0, nil)
	eventBus.On("Publish", mock.Anything, mock.Anything).Return(nil)

	fraudEngine := rules.NewEngine(
		rules.NewAmountRule(10000),
		rules.NewLocationRule(nil),
		rules.NewVelocityRule(10, txRepo),
	)

	handler := cmdhandler.NewAnalyzeTransactionHandler(txRepo, accRepo, eventBus, fraudEngine, log)

	cmd := cmdmodel.AnalyzeTransaction{
		AccountID:  uuid.New(),
		Amount:     100,
		Currency:   "BRL",
		MerchantID: "merchant-1",
		Location:   "BR",
	}

	result, err := handler.Handle(context.Background(), cmd)

	assert.NoError(t, err)
	assert.NotNil(t, result)

	txResult, ok := result.(cmdmodel.AnalyzeTransactionResult)
	assert.True(t, ok)
	assert.NotEqual(t, uuid.Nil, txResult.TransactionID)
	assert.Equal(t, "approved", txResult.Status)
	assert.Equal(t, "low", txResult.RiskLevel)

	txRepo.AssertExpectations(t)
	eventBus.AssertExpectations(t)
}

func TestAnalyzeTransactionHandler_HighRisk(t *testing.T) {
	txRepo := &mocks.TransactionWriteRepositoryMock{}
	accRepo := &mocks.AccountWriteRepositoryMock{}
	eventBus := &mocks.EventBusMock{}
	log := logger.NewNop()

	txRepo.On("Save", mock.Anything, mock.Anything).Return(nil)
	txRepo.On("CountRecentByAccountID", mock.Anything, mock.Anything, mock.Anything).Return(0, nil)
	eventBus.On("Publish", mock.Anything, mock.Anything).Return(nil)

	fraudEngine := rules.NewEngine(
		rules.NewAmountRule(100),   // triggers at 99999
		rules.NewLocationRule(nil), // triggers at "XX"
	)

	handler := cmdhandler.NewAnalyzeTransactionHandler(txRepo, accRepo, eventBus, fraudEngine, log)

	cmd := cmdmodel.AnalyzeTransaction{
		AccountID:  uuid.New(),
		Amount:     99999,
		Currency:   "BRL",
		MerchantID: "merchant-1",
		Location:   "XX",
	}

	result, err := handler.Handle(context.Background(), cmd)

	assert.NoError(t, err)
	txResult := result.(cmdmodel.AnalyzeTransactionResult)
	assert.Equal(t, "declined", txResult.Status)
	assert.Equal(t, "high", txResult.RiskLevel)
}

func TestAnalyzeTransactionHandler_InvalidCommand(t *testing.T) {
	handler := cmdhandler.NewAnalyzeTransactionHandler(nil, nil, nil, rules.NewEngine(), logger.NewNop())

	_, err := handler.Handle(context.Background(), "not-a-command")

	assert.Error(t, err)
}
