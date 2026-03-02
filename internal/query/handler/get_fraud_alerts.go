package handler

import (
	"context"
	"fmt"

	"github.com/hebertzin/cqrs/internal/application/bus"
	"github.com/hebertzin/cqrs/internal/query/repository"
)

const GetFraudAlertsKey = "GetFraudAlerts"

type GetFraudAlertsQuery struct {
	Page  int
	Limit int
}

type GetFraudAlertsHandler struct {
	readRepo repository.TransactionReadRepository
}

func NewGetFraudAlertsHandler(readRepo repository.TransactionReadRepository) *GetFraudAlertsHandler {
	return &GetFraudAlertsHandler{readRepo: readRepo}
}

func (h *GetFraudAlertsHandler) Handle(ctx context.Context, q bus.Query) (interface{}, error) {
	query, ok := q.(GetFraudAlertsQuery)
	if !ok {
		return nil, fmt.Errorf("invalid query type for GetFraudAlertsHandler")
	}

	if query.Page < 1 {
		query.Page = 1
	}
	if query.Limit < 1 || query.Limit > 100 {
		query.Limit = 20
	}

	result, err := h.readRepo.GetFraudAlerts(ctx, query.Page, query.Limit)
	if err != nil {
		return nil, fmt.Errorf("fetching fraud alerts: %w", err)
	}

	return result, nil
}
