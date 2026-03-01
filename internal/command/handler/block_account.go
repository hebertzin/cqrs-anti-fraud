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

type BlockAccountHandler struct {
	accountRepo repository.AccountWriteRepository
	eventBus    eventbus.EventBus
	logger      *zap.Logger
}

func NewBlockAccountHandler(
	accountRepo repository.AccountWriteRepository,
	eventBus eventbus.EventBus,
	logger *zap.Logger,
) *BlockAccountHandler {
	return &BlockAccountHandler{
		accountRepo: accountRepo,
		eventBus:    eventBus,
		logger:      logger,
	}
}

func (h *BlockAccountHandler) Handle(ctx context.Context, cmd bus.Command) (bus.CommandResult, error) {
	command, ok := cmd.(cmdmodel.BlockAccount)
	if !ok {
		return nil, fmt.Errorf("invalid command type for BlockAccountHandler")
	}

	account, err := h.accountRepo.FindByID(ctx, command.AccountID)
	if err != nil {
		return nil, fmt.Errorf("finding account %s: %w", command.AccountID, err)
	}

	if account.IsBlocked() {
		return nil, fmt.Errorf("account %s is already blocked", command.AccountID)
	}

	account.Block()

	if err := h.accountRepo.Update(ctx, account); err != nil {
		return nil, fmt.Errorf("updating account: %w", err)
	}

	h.logger.Info("account blocked",
		zap.String("account_id", command.AccountID.String()),
		zap.String("reason", command.Reason),
		zap.String("blocked_by", command.BlockedBy),
	)

	evt := event.NewAccountBlocked(command.AccountID, command.Reason, command.BlockedBy)
	if err := h.eventBus.Publish(ctx, evt); err != nil {
		h.logger.Error("failed to publish AccountBlocked event", zap.Error(err))
	}

	return nil, nil
}
