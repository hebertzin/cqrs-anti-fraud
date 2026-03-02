package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"github.com/hebertzin/cqrs/internal/query/model"
)

const (
	accountKeyPrefix = "account:"
	accountTTL       = 7 * 24 * time.Hour
)

type AccountRepository struct {
	client *redis.Client
}

func NewAccountRepository(client *redis.Client) *AccountRepository {
	return &AccountRepository{client: client}
}

func (r *AccountRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.AccountStatusView, error) {
	key := accountKeyPrefix + id.String()
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, fmt.Errorf("getting account view %s: %w", id, err)
	}

	var view model.AccountStatusView
	if err := json.Unmarshal(data, &view); err != nil {
		return nil, fmt.Errorf("unmarshaling account view: %w", err)
	}
	return &view, nil
}

func (r *AccountRepository) Save(ctx context.Context, view *model.AccountStatusView) error {
	data, err := json.Marshal(view)
	if err != nil {
		return fmt.Errorf("marshaling account view: %w", err)
	}

	key := accountKeyPrefix + view.ID.String()
	if err := r.client.Set(ctx, key, data, accountTTL).Err(); err != nil {
		return fmt.Errorf("saving account view: %w", err)
	}
	return nil
}

func (r *AccountRepository) IncrementTransactionCount(ctx context.Context, accountID uuid.UUID) error {
	return r.updateAccountView(ctx, accountID, func(v *model.AccountStatusView) {
		v.TotalTransactions++
		now := time.Now().UTC()
		v.LastActivityAt = &now
	})
}

func (r *AccountRepository) IncrementFlaggedCount(ctx context.Context, accountID uuid.UUID) error {
	return r.updateAccountView(ctx, accountID, func(v *model.AccountStatusView) {
		v.FlaggedCount++
	})
}

func (r *AccountRepository) IncrementDeclinedCount(ctx context.Context, accountID uuid.UUID) error {
	return r.updateAccountView(ctx, accountID, func(v *model.AccountStatusView) {
		v.DeclinedCount++
	})
}

func (r *AccountRepository) updateAccountView(
	ctx context.Context, accountID uuid.UUID, update func(*model.AccountStatusView),
) error {
	view, err := r.GetByID(ctx, accountID)
	if err != nil {
		// If view doesn't exist, create a new one
		view = &model.AccountStatusView{
			ID:        accountID,
			Status:    "active",
			CreatedAt: time.Now().UTC(),
		}
	}

	update(view)
	return r.Save(ctx, view)
}
