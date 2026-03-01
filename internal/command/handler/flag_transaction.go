package handler

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/hebertzin/cqrs/internal/application/bus"
	"github.com/hebertzin/cqrs/internal/application/eventbus"
	cmdmodel "github.com/hebertzin/cqrs/internal/command/model"
	"github.com/hebertzin/cqrs/internal/domain/event"
	"github.com/hebertzin/cqrs/internal/domain/repository"
)

type FlagTransactionHandler struct {
	writeRepo repository.TransactionWriteRepository
	eventBus  eventbus.EventBus
	logger    *zap.Logger
}

func NewFlagTransactionHandler(
	writeRepo repository.TransactionWriteRepository,
	eventBus eventbus.EventBus,
	logger *zap.Logger,
) *FlagTransactionHandler {
	return &FlagTransactionHandler{
		writeRepo: writeRepo,
		eventBus:  eventBus,
		logger:    logger,
	}
}

func (h *FlagTransactionHandler) Handle(ctx context.Context, cmd bus.Command) (bus.CommandResult, error) {
	command, ok := cmd.(cmdmodel.FlagTransaction)
	if !ok {
		return nil, fmt.Errorf("invalid command type for FlagTransactionHandler")
	}

	tx, err := h.writeRepo.FindByID(ctx, command.TransactionID)
	if err != nil {
		return nil, fmt.Errorf("finding transaction %s: %w", command.TransactionID, err)
	}

	tx.Flag()

	if err := h.writeRepo.Update(ctx, tx); err != nil {
		return nil, fmt.Errorf("updating transaction: %w", err)
	}

	h.logger.Info("transaction flagged manually",
		zap.String("transaction_id", command.TransactionID.String()),
		zap.String("reason", command.Reason),
		zap.String("flagged_by", command.FlaggedBy),
	)

	evt := event.NewTransactionFlagged(command.TransactionID, tx.AccountID, command.Reason, command.FlaggedBy)
	if err := h.eventBus.Publish(ctx, evt); err != nil {
		h.logger.Error("failed to publish TransactionFlagged event", zap.Error(err))
	}

	return nil, nil
}
