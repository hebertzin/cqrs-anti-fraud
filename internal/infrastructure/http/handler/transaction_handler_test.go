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

type fakeFlagTxHandler struct{}

func (h *fakeFlagTxHandler) Handle(_ context.Context, _ bus.Command) (bus.CommandResult, error) {
	return nil, nil
}

// fakeAnalyzeTxHandler is a test double for the analyze transaction command handler.
type fakeAnalyzeTxHandler struct{}

func (h *fakeAnalyzeTxHandler) Handle(_ context.Context, _ bus.Command) (bus.CommandResult, error) {
	return cmdmodel.AnalyzeTransactionResult{
		TransactionID: uuid.New(),
		Status:        "approved",
		RiskScore:     0.1,
		RiskLevel:     "low",
	}, nil
}

type fakeGetTxRiskHandler struct {
	view *model.TransactionRiskView
}

func (h *fakeGetTxRiskHandler) Handle(_ context.Context, _ bus.Query) (interface{}, error) {
	return h.view, nil
}

func buildTransactionHandler(t *testing.T) (*handler.TransactionHandler, *chi.Mux) {
	t.Helper()
	commandBus := bus.NewCommandBus()
	commandBus.Register(cmdmodel.AnalyzeTransactionCommand, &fakeAnalyzeTxHandler{})
	commandBus.Register(cmdmodel.FlagTransactionCommand, &fakeFlagTxHandler{})

	txID := uuid.New()
	view := &model.TransactionRiskView{
		ID:        txID,
		AccountID: uuid.New(),
		Amount:    100,
		Currency:  "BRL",
		RiskScore: 0.1,
		RiskLevel: model.RiskLevelLow,
		Status:    "approved",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	queryBus := bus.NewQueryBus()
	queryBus.Register(queryhandler.GetTransactionRiskKey, &fakeGetTxRiskHandler{view: view})

	h := handler.NewTransactionHandler(commandBus, queryBus, zap.NewNop())

	r := chi.NewRouter()
	r.Post("/transactions", h.AnalyzeTransaction)
	r.Get("/transactions/{id}/risk", h.GetTransactionRisk)
	r.Post("/transactions/{id}/flag", h.FlagTransaction)

	return h, r
}

func TestTransactionHandler_AnalyzeTransaction_Success(t *testing.T) {
	_, r := buildTransactionHandler(t)

	body, _ := json.Marshal(map[string]interface{}{
		"account_id":  uuid.New().String(),
		"amount":      500.0,
		"currency":    "BRL",
		"merchant_id": "merchant-1",
		"location":    "BR",
	})

	req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "approved", resp["status"])
}

func TestTransactionHandler_AnalyzeTransaction_InvalidBody(t *testing.T) {
	_, r := buildTransactionHandler(t)

	req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewReader([]byte("not-json")))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestTransactionHandler_AnalyzeTransaction_InvalidAccountID(t *testing.T) {
	_, r := buildTransactionHandler(t)

	body, _ := json.Marshal(map[string]interface{}{
		"account_id": "not-a-uuid",
		"amount":     100.0,
		"currency":   "BRL",
	})
	req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestTransactionHandler_AnalyzeTransaction_ZeroAmount(t *testing.T) {
	_, r := buildTransactionHandler(t)

	body, _ := json.Marshal(map[string]interface{}{
		"account_id": uuid.New().String(),
		"amount":     0,
		"currency":   "BRL",
	})
	req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestTransactionHandler_AnalyzeTransaction_InvalidCurrency(t *testing.T) {
	_, r := buildTransactionHandler(t)

	body, _ := json.Marshal(map[string]interface{}{
		"account_id": uuid.New().String(),
		"amount":     100.0,
		"currency":   "BR",
	})
	req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestTransactionHandler_GetTransactionRisk_Success(t *testing.T) {
	_, r := buildTransactionHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/transactions/"+uuid.New().String()+"/risk", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestTransactionHandler_GetTransactionRisk_InvalidID(t *testing.T) {
	_, r := buildTransactionHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/transactions/not-a-uuid/risk", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestTransactionHandler_FlagTransaction_Success(t *testing.T) {
	_, r := buildTransactionHandler(t)

	body, _ := json.Marshal(map[string]interface{}{
		"reason":     "suspicious activity",
		"flagged_by": "analyst",
	})
	req := httptest.NewRequest(http.MethodPost, "/transactions/"+uuid.New().String()+"/flag", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestTransactionHandler_FlagTransaction_InvalidID(t *testing.T) {
	_, r := buildTransactionHandler(t)

	req := httptest.NewRequest(http.MethodPost, "/transactions/not-a-uuid/flag", bytes.NewReader([]byte("{}")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestTransactionHandler_FlagTransaction_InvalidBody(t *testing.T) {
	_, r := buildTransactionHandler(t)

	req := httptest.NewRequest(http.MethodPost, "/transactions/"+uuid.New().String()+"/flag", bytes.NewReader([]byte("not-json")))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
