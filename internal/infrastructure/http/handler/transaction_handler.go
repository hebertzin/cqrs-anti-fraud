package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/hebertzin/cqrs/internal/application/bus"
	cmdmodel "github.com/hebertzin/cqrs/internal/command/model"
	queryhandler "github.com/hebertzin/cqrs/internal/query/handler"
	"github.com/hebertzin/cqrs/internal/infrastructure/http/response"
)

type TransactionHandler struct {
	commandBus *bus.CommandBus
	queryBus   *bus.QueryBus
	logger     *zap.Logger
}

func NewTransactionHandler(commandBus *bus.CommandBus, queryBus *bus.QueryBus, logger *zap.Logger) *TransactionHandler {
	return &TransactionHandler{
		commandBus: commandBus,
		queryBus:   queryBus,
		logger:     logger,
	}
}

type analyzeTransactionRequest struct {
	AccountID  string            `json:"account_id"`
	Amount     float64           `json:"amount"`
	Currency   string            `json:"currency"`
	MerchantID string            `json:"merchant_id"`
	Location   string            `json:"location"`
	Metadata   map[string]string `json:"metadata,omitempty"`
}

func (h *TransactionHandler) AnalyzeTransaction(w http.ResponseWriter, r *http.Request) {
	var req analyzeTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	accountID, err := uuid.Parse(req.AccountID)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid account_id format")
		return
	}

	if req.Amount <= 0 {
		response.Error(w, http.StatusBadRequest, "amount must be greater than zero")
		return
	}
	if len(req.Currency) != 3 {
		response.Error(w, http.StatusBadRequest, "currency must be a 3-letter ISO code")
		return
	}

	cmd := cmdmodel.AnalyzeTransaction{
		AccountID:  accountID,
		Amount:     req.Amount,
		Currency:   req.Currency,
		MerchantID: req.MerchantID,
		Location:   req.Location,
		Metadata:   req.Metadata,
	}

	result, err := h.commandBus.Dispatch(r.Context(), cmdmodel.AnalyzeTransactionCommand, cmd)
	if err != nil {
		h.logger.Error("failed to analyze transaction", zap.Error(err))
		response.Error(w, http.StatusInternalServerError, "failed to analyze transaction")
		return
	}

	response.Created(w, result)
}

func (h *TransactionHandler) GetTransactionRisk(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid transaction id")
		return
	}

	result, err := h.queryBus.Query(r.Context(), queryhandler.GetTransactionRiskKey, queryhandler.GetTransactionRiskQuery{TransactionID: id})
	if err != nil {
		h.logger.Error("failed to get transaction risk", zap.Error(err))
		response.Error(w, http.StatusNotFound, "transaction not found")
		return
	}

	response.OK(w, result)
}

func (h *TransactionHandler) FlagTransaction(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid transaction id")
		return
	}

	var req struct {
		Reason    string `json:"reason"`
		FlaggedBy string `json:"flagged_by"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	cmd := cmdmodel.FlagTransaction{
		TransactionID: id,
		Reason:        req.Reason,
		FlaggedBy:     req.FlaggedBy,
	}

	if _, err := h.commandBus.Dispatch(r.Context(), cmdmodel.FlagTransactionCommand, cmd); err != nil {
		h.logger.Error("failed to flag transaction", zap.Error(err))
		response.Error(w, http.StatusInternalServerError, "failed to flag transaction")
		return
	}

	response.NoContent(w)
}
