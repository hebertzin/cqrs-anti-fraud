package handler

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/hebertzin/cqrs/internal/application/bus"
	"github.com/hebertzin/cqrs/internal/application/eventbus"
	cmdmodel "github.com/hebertzin/cqrs/internal/command/model"
	"github.com/hebertzin/cqrs/internal/domain/entity"
	"github.com/hebertzin/cqrs/internal/domain/event"
	"github.com/hebertzin/cqrs/internal/domain/repository"
	"github.com/hebertzin/cqrs/internal/fraud/rules"
	qmodel "github.com/hebertzin/cqrs/internal/query/model"
)

type AnalyzeTransactionHandler struct {
	writeRepo   repository.TransactionWriteRepository
	accountRepo repository.AccountWriteRepository
	eventBus    eventbus.EventBus
	fraudEngine *rules.Engine
	logger      *zap.Logger
}

func NewAnalyzeTransactionHandler(
	writeRepo repository.TransactionWriteRepository,
	accountRepo repository.AccountWriteRepository,
	eventBus eventbus.EventBus,
	fraudEngine *rules.Engine,
	logger *zap.Logger,
) *AnalyzeTransactionHandler {
	return &AnalyzeTransactionHandler{
		writeRepo:   writeRepo,
		accountRepo: accountRepo,
		eventBus:    eventBus,
		fraudEngine: fraudEngine,
		logger:      logger,
	}
}

func (h *AnalyzeTransactionHandler) Handle(ctx context.Context, cmd bus.Command) (bus.CommandResult, error) {
	command, ok := cmd.(cmdmodel.AnalyzeTransaction)
	if !ok {
		return nil, fmt.Errorf("invalid command type for AnalyzeTransactionHandler")
	}

	tx := entity.NewTransaction(
		command.AccountID,
		command.Amount,
		command.Currency,
		command.MerchantID,
		command.Location,
	)
	tx.Metadata = command.Metadata

	evaluation := h.fraudEngine.Evaluate(ctx, tx)
	tx.SetRiskScore(evaluation.TotalScore)

	switch {
	case tx.IsHighRisk():
		tx.Decline()
	case tx.IsMediumRisk():
		tx.Flag()
	default:
		tx.Approve()
	}

	if err := h.writeRepo.Save(ctx, tx); err != nil {
		return nil, fmt.Errorf("saving transaction: %w", err)
	}

	h.logger.Info("transaction analyzed",
		zap.String("transaction_id", tx.ID.String()),
		zap.String("status", string(tx.Status)),
		zap.Float64("risk_score", tx.RiskScore),
	)

	evt := event.NewTransactionAnalyzed(tx, evaluation.Reasons)
	if err := h.eventBus.Publish(ctx, evt); err != nil {
		h.logger.Error("failed to publish TransactionAnalyzed event", zap.Error(err))
	}

	return cmdmodel.AnalyzeTransactionResult{
		TransactionID: tx.ID,
		Status:        string(tx.Status),
		RiskScore:     tx.RiskScore,
		RiskLevel:     string(qmodel.RiskLevelFromScore(tx.RiskScore)),
	}, nil
}
