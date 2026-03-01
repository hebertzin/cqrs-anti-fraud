package projector

import (
	"context"

	"go.uber.org/zap"

	"github.com/hebertzin/cqrs/internal/application/eventbus"
	"github.com/hebertzin/cqrs/internal/domain/event"
	"github.com/hebertzin/cqrs/internal/query/repository"
)

// AccountProjector listens to account events and updates the read model.
type AccountProjector struct {
	accountReadRepo repository.AccountReadRepository
	logger          *zap.Logger
}

func NewAccountProjector(accountReadRepo repository.AccountReadRepository, logger *zap.Logger) *AccountProjector {
	return &AccountProjector{accountReadRepo: accountReadRepo, logger: logger}
}

func (p *AccountProjector) Register(bus eventbus.EventBus) {
	bus.Subscribe(event.TypeAccountBlocked, p.onAccountBlocked)
}

func (p *AccountProjector) onAccountBlocked(ctx context.Context, e event.Event) error {
	evt, ok := e.(event.AccountBlocked)
	if !ok {
		return nil
	}

	view, err := p.accountReadRepo.GetByID(ctx, evt.AccountID)
	if err != nil {
		p.logger.Warn("account view not found for blocking", zap.String("id", evt.AccountID.String()))
		return nil
	}

	now := evt.OccurredAt
	view.Status = "blocked"
	view.BlockedAt = &now

	return p.accountReadRepo.Save(ctx, view)
}
