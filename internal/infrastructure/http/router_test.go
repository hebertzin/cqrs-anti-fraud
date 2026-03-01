package http_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/hebertzin/cqrs/internal/application/bus"
	cmdmodel "github.com/hebertzin/cqrs/internal/command/model"
	infrahttp "github.com/hebertzin/cqrs/internal/infrastructure/http"
	"github.com/hebertzin/cqrs/internal/infrastructure/http/handler"
	queryhandler "github.com/hebertzin/cqrs/internal/query/handler"
	"github.com/hebertzin/cqrs/internal/query/model"
)

// stub handlers for wiring the router in tests

type noopCmdHandler struct{}

func (h *noopCmdHandler) Handle(_ context.Context, _ bus.Command) (bus.CommandResult, error) {
	return nil, nil
}

type noopQueryHandler struct{}

func (h *noopQueryHandler) Handle(_ context.Context, _ bus.Query) (interface{}, error) {
	return &model.TransactionRiskView{}, nil
}

type noopAccountQueryHandler struct{}

func (h *noopAccountQueryHandler) Handle(_ context.Context, _ bus.Query) (interface{}, error) {
	return &model.AccountStatusView{}, nil
}

type noopFraudAlertsHandler struct{}

func (h *noopFraudAlertsHandler) Handle(_ context.Context, _ bus.Query) (interface{}, error) {
	return &model.FraudAlertListResponse{}, nil
}

func buildRouter() http.Handler {
	log := zap.NewNop()

	cmdBus := bus.NewCommandBus()
	cmdBus.Register(cmdmodel.AnalyzeTransactionCommand, &noopCmdHandler{})
	cmdBus.Register(cmdmodel.BlockAccountCommand, &noopCmdHandler{})
	cmdBus.Register(cmdmodel.FlagTransactionCommand, &noopCmdHandler{})

	qBus := bus.NewQueryBus()
	qBus.Register(queryhandler.GetTransactionRiskKey, &noopQueryHandler{})
	qBus.Register(queryhandler.GetAccountStatusKey, &noopAccountQueryHandler{})
	qBus.Register(queryhandler.GetFraudAlertsKey, &noopFraudAlertsHandler{})

	txHandler := handler.NewTransactionHandler(cmdBus, qBus, log)
	accHandler := handler.NewAccountHandler(cmdBus, qBus, log)
	healthHandler := handler.NewHealthHandler()

	return infrahttp.NewRouter(txHandler, accHandler, healthHandler, log)
}

func TestRouter_HealthEndpoint(t *testing.T) {
	r := buildRouter()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRouter_ReadyEndpoint(t *testing.T) {
	r := buildRouter()
	req := httptest.NewRequest(http.MethodGet, "/ready", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRouter_UnknownRoute(t *testing.T) {
	r := buildRouter()
	req := httptest.NewRequest(http.MethodGet, "/unknown", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouter_MethodNotAllowed(t *testing.T) {
	r := buildRouter()
	req := httptest.NewRequest(http.MethodDelete, "/health", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}
