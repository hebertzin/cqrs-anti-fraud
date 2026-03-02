package handler_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	queryhandler "github.com/hebertzin/cqrs/internal/query/handler"
	"github.com/hebertzin/cqrs/internal/query/model"
	"github.com/hebertzin/cqrs/tests/mocks"
)

func TestGetFraudAlertsHandler_Success(t *testing.T) {
	readRepo := &mocks.TransactionReadRepositoryMock{}
	expected := &model.FraudAlertListResponse{
		Alerts: []*model.FraudAlertView{},
		Total:  0,
		Page:   1,
		Limit:  20,
	}
	readRepo.On("GetFraudAlerts", mock.Anything, 1, 20).Return(expected, nil)

	handler := queryhandler.NewGetFraudAlertsHandler(readRepo)
	result, err := handler.Handle(context.Background(), queryhandler.GetFraudAlertsQuery{Page: 1, Limit: 20})

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestGetFraudAlertsHandler_DefaultPagination(t *testing.T) {
	readRepo := &mocks.TransactionReadRepositoryMock{}
	readRepo.On("GetFraudAlerts", mock.Anything, 1, 20).Return(&model.FraudAlertListResponse{}, nil)

	handler := queryhandler.NewGetFraudAlertsHandler(readRepo)
	// Zero values should be normalized to defaults
	_, err := handler.Handle(context.Background(), queryhandler.GetFraudAlertsQuery{Page: 0, Limit: 0})

	assert.NoError(t, err)
	readRepo.AssertExpectations(t)
}

func TestGetFraudAlertsHandler_Error(t *testing.T) {
	readRepo := &mocks.TransactionReadRepositoryMock{}
	readRepo.On("GetFraudAlerts", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("redis error"))

	handler := queryhandler.NewGetFraudAlertsHandler(readRepo)
	_, err := handler.Handle(context.Background(), queryhandler.GetFraudAlertsQuery{Page: 1, Limit: 20})

	assert.Error(t, err)
}

func TestGetFraudAlertsHandler_InvalidQuery(t *testing.T) {
	handler := queryhandler.NewGetFraudAlertsHandler(&mocks.TransactionReadRepositoryMock{})
	_, err := handler.Handle(context.Background(), "not-a-query")
	assert.Error(t, err)
}
