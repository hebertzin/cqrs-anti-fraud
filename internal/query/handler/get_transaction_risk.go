package handler

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/hebertzin/cqrs/internal/application/bus"
	"github.com/hebertzin/cqrs/internal/query/model"
	"github.com/hebertzin/cqrs/internal/query/repository"
)

const GetTransactionRiskKey = "GetTransactionRisk"

type GetTransactionRiskQuery struct {
	TransactionID uuid.UUID
}

type GetTransactionRiskHandler struct {
	readRepo repository.TransactionReadRepository
}

func NewGetTransactionRiskHandler(readRepo repository.TransactionReadRepository) *GetTransactionRiskHandler {
	return &GetTransactionRiskHandler{readRepo: readRepo}
}

func (h *GetTransactionRiskHandler) Handle(ctx context.Context, q bus.Query) (interface{}, error) {
	query, ok := q.(GetTransactionRiskQuery)
	if !ok {
		return nil, fmt.Errorf("invalid query type for GetTransactionRiskHandler")
	}

	view, err := h.readRepo.GetByID(ctx, query.TransactionID)
	if err != nil {
		return nil, fmt.Errorf("fetching transaction risk view: %w", err)
	}

	view.RiskLevel = model.RiskLevelFromScore(view.RiskScore)
	return view, nil
}
