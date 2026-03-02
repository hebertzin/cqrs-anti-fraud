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
	txKeyPrefix   = "tx:"
	fraudAlertKey = "fraud:alerts"
	txTTL         = 24 * time.Hour
)

type TransactionRepository struct {
	client *redis.Client
}

func NewTransactionRepository(client *redis.Client) *TransactionRepository {
	return &TransactionRepository{client: client}
}

func (r *TransactionRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.TransactionRiskView, error) {
	key := txKeyPrefix + id.String()
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, fmt.Errorf("getting transaction view %s: %w", id, err)
	}

	var view model.TransactionRiskView
	if err := json.Unmarshal(data, &view); err != nil {
		return nil, fmt.Errorf("unmarshaling transaction view: %w", err)
	}
	return &view, nil
}

func (r *TransactionRepository) GetByAccountID(ctx context.Context, accountID uuid.UUID) ([]*model.TransactionRiskView, error) {
	pattern := txKeyPrefix + "account:" + accountID.String() + ":*"
	keys, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, fmt.Errorf("listing transaction keys: %w", err)
	}

	var views []*model.TransactionRiskView
	for _, key := range keys {
		data, err := r.client.Get(ctx, key).Bytes()
		if err != nil {
			continue
		}
		var view model.TransactionRiskView
		if err := json.Unmarshal(data, &view); err != nil {
			continue
		}
		views = append(views, &view)
	}
	return views, nil
}

func (r *TransactionRepository) Save(ctx context.Context, view *model.TransactionRiskView) error {
	data, err := json.Marshal(view)
	if err != nil {
		return fmt.Errorf("marshaling transaction view: %w", err)
	}

	key := txKeyPrefix + view.ID.String()
	if err := r.client.Set(ctx, key, data, txTTL).Err(); err != nil {
		return fmt.Errorf("saving transaction view: %w", err)
	}

	if view.RiskScore >= 0.5 {
		alertData, err := json.Marshal(buildFraudAlert(view))
		if err == nil {
			r.client.LPush(ctx, fraudAlertKey, alertData)
			r.client.LTrim(ctx, fraudAlertKey, 0, 9999) // keep last 10k alerts
		}
	}

	return nil
}

func (r *TransactionRepository) GetFraudAlerts(ctx context.Context, page, limit int) (*model.FraudAlertListResponse, error) {
	start := int64((page - 1) * limit)
	stop := start + int64(limit) - 1

	items, err := r.client.LRange(ctx, fraudAlertKey, start, stop).Result()
	if err != nil {
		return nil, fmt.Errorf("listing fraud alerts: %w", err)
	}

	total, _ := r.client.LLen(ctx, fraudAlertKey).Result()

	var alerts []*model.FraudAlertView
	for _, item := range items {
		var alert model.FraudAlertView
		if err := json.Unmarshal([]byte(item), &alert); err == nil {
			alerts = append(alerts, &alert)
		}
	}

	return &model.FraudAlertListResponse{
		Alerts: alerts,
		Total:  int(total),
		Page:   page,
		Limit:  limit,
	}, nil
}

func buildFraudAlert(view *model.TransactionRiskView) model.FraudAlertView {
	return model.FraudAlertView{
		ID:            uuid.New(),
		TransactionID: view.ID,
		AccountID:     view.AccountID,
		Amount:        view.Amount,
		Currency:      view.Currency,
		RiskScore:     view.RiskScore,
		RiskLevel:     model.RiskLevelFromScore(view.RiskScore),
		Reasons:       view.FraudReasons,
		Status:        model.FraudAlertStatusOpen,
		CreatedAt:     time.Now().UTC(),
	}
}
