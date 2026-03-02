package handler

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hebertzin/cqrs/internal/application/bus"
	"github.com/hebertzin/cqrs/internal/query/repository"
)

const GetAccountStatusKey = "GetAccountStatus"

type GetAccountStatusQuery struct {
	AccountID uuid.UUID
}

type GetAccountStatusHandler struct {
	readRepo repository.AccountReadRepository
}

func NewGetAccountStatusHandler(readRepo repository.AccountReadRepository) *GetAccountStatusHandler {
	return &GetAccountStatusHandler{readRepo: readRepo}
}

func (h *GetAccountStatusHandler) Handle(ctx context.Context, q bus.Query) (interface{}, error) {
	query, ok := q.(GetAccountStatusQuery)
	if !ok {
		return nil, fmt.Errorf("invalid query type for GetAccountStatusHandler")
	}

	view, err := h.readRepo.GetByID(ctx, query.AccountID)
	if err != nil {
		return nil, fmt.Errorf("fetching account status view: %w", err)
	}

	return view, nil
}
