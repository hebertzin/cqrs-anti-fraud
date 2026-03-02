package projector

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/hebertzin/cqrs/internal/application/eventbus"
	"github.com/hebertzin/cqrs/internal/domain/event"
	"github.com/hebertzin/cqrs/internal/query/model"
	"github.com/hebertzin/cqrs/internal/query/repository"
)

type TransactionProjector struct {
	txReadRepo      repository.TransactionReadRepository
	accountReadRepo repository.AccountReadRepository
	logger          *zap.Logger
}

func NewTransactionProjector(
	txReadRepo repository.TransactionReadRepository,
	accountReadRepo repository.AccountReadRepository,
	logger *zap.Logger,
) *TransactionProjector {
	return &TransactionProjector{
		txReadRepo:      txReadRepo,
		accountReadRepo: accountReadRepo,
		logger:          logger,
	}
}

func (p *TransactionProjector) Register(bus eventbus.EventBus) {
	bus.Subscribe(event.TypeTransactionAnalyzed, p.onTransactionAnalyzed)
	bus.Subscribe(event.TypeTransactionFlagged, p.onTransactionFlagged)
}

func (p *TransactionProjector) onTransactionAnalyzed(ctx context.Context, e event.Event) error {
	evt, ok := e.(event.TransactionAnalyzed)
	if !ok {
		return nil
	}

	view := &model.TransactionRiskView{
		ID:           evt.TransactionID,
		AccountID:    evt.AccountID,
		Amount:       evt.Amount,
		Currency:     evt.Currency,
		MerchantID:   evt.MerchantID,
		Location:     evt.Location,
		Status:       string(evt.Status),
		RiskScore:    evt.RiskScore,
		RiskLevel:    model.RiskLevelFromScore(evt.RiskScore),
		FraudReasons: evt.FraudReasons,
		CreatedAt:    evt.OccurredAt,
		UpdatedAt:    evt.OccurredAt,
	}

	if err := p.txReadRepo.Save(ctx, view); err != nil {
		p.logger.Error("failed to project transaction view", zap.Error(err))
		return fmt.Errorf("project transaction view: %w", err)
	}

	if err := p.accountReadRepo.IncrementTransactionCount(ctx, evt.AccountID); err != nil {
		p.logger.Error("failed to increment transaction count", zap.Error(err))
	}

	if evt.RiskScore >= 0.5 {
		if err := p.accountReadRepo.IncrementFlaggedCount(ctx, evt.AccountID); err != nil {
			p.logger.Error("failed to increment flagged count", zap.Error(err))
		}
	}

	if evt.RiskScore >= 0.8 {
		if err := p.accountReadRepo.IncrementDeclinedCount(ctx, evt.AccountID); err != nil {
			p.logger.Error("failed to increment declined count", zap.Error(err))
		}
	}

	return nil
}

func (p *TransactionProjector) onTransactionFlagged(ctx context.Context, e event.Event) error {
	evt, ok := e.(event.TransactionFlagged)
	if !ok {
		return nil
	}

	view, err := p.txReadRepo.GetByID(ctx, evt.TransactionID)
	if err != nil {
		p.logger.Warn("transaction view not found for flagging", zap.String("id", evt.TransactionID.String()))
		return nil
	}

	view.Status = "flagged"
	view.UpdatedAt = evt.OccurredAt

	if err := p.txReadRepo.Save(ctx, view); err != nil {
		return fmt.Errorf("update flagged transaction view: %w", err)
	}
	return nil
}
