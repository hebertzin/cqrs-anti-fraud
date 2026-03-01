package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hebertzin/cqrs/internal/domain/entity"
)

type AccountRepository struct {
	pool *pgxpool.Pool
}

func NewAccountRepository(pool *pgxpool.Pool) *AccountRepository {
	return &AccountRepository{pool: pool}
}

func (r *AccountRepository) Save(ctx context.Context, account *entity.Account) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO accounts (id, user_id, status, risk_level, blocked_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		account.ID, account.UserID, string(account.Status), account.RiskLevel,
		account.BlockedAt, account.CreatedAt, account.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("inserting account: %w", err)
	}
	return nil
}

func (r *AccountRepository) Update(ctx context.Context, account *entity.Account) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE accounts SET status = $1, risk_level = $2, blocked_at = $3, updated_at = $4
		WHERE id = $5`,
		string(account.Status), account.RiskLevel, account.BlockedAt, account.UpdatedAt, account.ID,
	)
	if err != nil {
		return fmt.Errorf("updating account: %w", err)
	}
	return nil
}

func (r *AccountRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Account, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, user_id, status, risk_level, blocked_at, created_at, updated_at
		FROM accounts WHERE id = $1`, id)

	var account entity.Account
	var status string

	err := row.Scan(
		&account.ID, &account.UserID, &status, &account.RiskLevel,
		&account.BlockedAt, &account.CreatedAt, &account.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("scanning account row: %w", err)
	}

	account.Status = entity.AccountStatus(status)
	return &account, nil
}
