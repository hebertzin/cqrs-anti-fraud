package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/hebertzin/cqrs/internal/application/bus"
	cmdmodel "github.com/hebertzin/cqrs/internal/command/model"
	"github.com/hebertzin/cqrs/internal/infrastructure/http/handler"
	queryhandler "github.com/hebertzin/cqrs/internal/query/handler"
	"github.com/hebertzin/cqrs/internal/query/model"
)

type fakeGetAccountStatusHandler struct {
	view *model.AccountStatusView
}

func (h *fakeGetAccountStatusHandler) Handle(_ context.Context, _ bus.Query) (interface{}, error) {
	return h.view, nil
}

type fakeBlockAccountHandler struct{}

func (h *fakeBlockAccountHandler) Handle(_ context.Context, _ bus.Command) (bus.CommandResult, error) {
	return nil, nil
}

type fakeGetFraudAlertsHandler struct{}

func (h *fakeGetFraudAlertsHandler) Handle(_ context.Context, _ bus.Query) (interface{}, error) {
	return &model.FraudAlertListResponse{Alerts: nil, Total: 0, Page: 1, Limit: 20}, nil
}

func buildAccountHandler(t *testing.T) (*handler.AccountHandler, *chi.Mux) {
	t.Helper()

	now := time.Now()
	accountID := uuid.New()
	view := &model.AccountStatusView{
		ID:        accountID,
		Status:    "active",
		CreatedAt: now,
	}

	commandBus := bus.NewCommandBus()
	commandBus.Register(cmdmodel.BlockAccountCommand, &fakeBlockAccountHandler{})

	queryBus := bus.NewQueryBus()
	queryBus.Register(queryhandler.GetAccountStatusKey, &fakeGetAccountStatusHandler{view: view})
	queryBus.Register(queryhandler.GetFraudAlertsKey, &fakeGetFraudAlertsHandler{})

	h := handler.NewAccountHandler(commandBus, queryBus, zap.NewNop())

	r := chi.NewRouter()
	r.Get("/accounts/{id}/status", h.GetAccountStatus)
	r.Post("/accounts/{id}/block", h.BlockAccount)
	r.Get("/accounts/fraud-alerts", h.GetFraudAlerts)

	return h, r
}

func TestAccountHandler_GetAccountStatus_Success(t *testing.T) {
	_, r := buildAccountHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/accounts/"+uuid.New().String()+"/status", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "active", resp["status"])
}

func TestAccountHandler_GetAccountStatus_InvalidID(t *testing.T) {
	_, r := buildAccountHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/accounts/not-a-uuid/status", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAccountHandler_BlockAccount_Success(t *testing.T) {
	_, r := buildAccountHandler(t)

	body, _ := json.Marshal(map[string]interface{}{
		"reason":     "fraud suspected",
		"blocked_by": "analyst",
	})
	req := httptest.NewRequest(http.MethodPost, "/accounts/"+uuid.New().String()+"/block", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestAccountHandler_BlockAccount_InvalidID(t *testing.T) {
	_, r := buildAccountHandler(t)

	req := httptest.NewRequest(http.MethodPost, "/accounts/not-a-uuid/block", bytes.NewReader([]byte("{}")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAccountHandler_BlockAccount_InvalidBody(t *testing.T) {
	_, r := buildAccountHandler(t)

	req := httptest.NewRequest(http.MethodPost, "/accounts/"+uuid.New().String()+"/block", bytes.NewReader([]byte("not-json")))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAccountHandler_GetFraudAlerts_Success(t *testing.T) {
	_, r := buildAccountHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/accounts/fraud-alerts", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
