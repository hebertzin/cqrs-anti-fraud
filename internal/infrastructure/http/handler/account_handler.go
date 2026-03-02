package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/hebertzin/cqrs/internal/application/bus"
	cmdmodel "github.com/hebertzin/cqrs/internal/command/model"
	"github.com/hebertzin/cqrs/internal/infrastructure/http/response"
	queryhandler "github.com/hebertzin/cqrs/internal/query/handler"
)

type AccountHandler struct {
	commandBus *bus.CommandBus
	queryBus   *bus.QueryBus
	logger     *zap.Logger
}

func NewAccountHandler(commandBus *bus.CommandBus, queryBus *bus.QueryBus, logger *zap.Logger) *AccountHandler {
	return &AccountHandler{
		commandBus: commandBus,
		queryBus:   queryBus,
		logger:     logger,
	}
}

func (h *AccountHandler) GetAccountStatus(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid account id")
		return
	}

	query := queryhandler.GetAccountStatusQuery{AccountID: id}
	result, err := h.queryBus.Query(r.Context(), queryhandler.GetAccountStatusKey, query)
	if err != nil {
		h.logger.Error("failed to get account status", zap.Error(err))
		response.Error(w, http.StatusNotFound, "account not found")
		return
	}

	response.OK(w, result)
}

func (h *AccountHandler) BlockAccount(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid account id")
		return
	}

	var req struct {
		Reason    string `json:"reason"`
		BlockedBy string `json:"blocked_by"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	cmd := cmdmodel.BlockAccount{
		AccountID: id,
		Reason:    req.Reason,
		BlockedBy: req.BlockedBy,
	}

	if _, err := h.commandBus.Dispatch(r.Context(), cmdmodel.BlockAccountCommand, cmd); err != nil {
		h.logger.Error("failed to block account", zap.Error(err))
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.NoContent(w)
}

func (h *AccountHandler) GetFraudAlerts(w http.ResponseWriter, r *http.Request) {
	page := 1
	limit := 20

	result, err := h.queryBus.Query(r.Context(), queryhandler.GetFraudAlertsKey, queryhandler.GetFraudAlertsQuery{
		Page:  page,
		Limit: limit,
	})
	if err != nil {
		h.logger.Error("failed to get fraud alerts", zap.Error(err))
		response.Error(w, http.StatusInternalServerError, "failed to get fraud alerts")
		return
	}

	response.OK(w, result)
}
