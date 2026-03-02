package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hebertzin/cqrs/internal/domain/entity"
)

type TransactionRepository struct {
	pool *pgxpool.Pool
}

func NewTransactionRepository(pool *pgxpool.Pool) *TransactionRepository {
	return &TransactionRepository{pool: pool}
}

func (r *TransactionRepository) Save(ctx context.Context, tx *entity.Transaction) error {
	metadata, err := json.Marshal(tx.Metadata)
	if err != nil {
		return fmt.Errorf("marshaling metadata: %w", err)
	}

	_, err = r.pool.Exec(ctx, `
		INSERT INTO transactions
			(id, account_id, amount, currency, merchant_id, location, status, risk_score, metadata, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		tx.ID, tx.AccountID, tx.Amount, tx.Currency, tx.MerchantID, tx.Location,
		string(tx.Status), tx.RiskScore, metadata, tx.CreatedAt, tx.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("inserting transaction: %w", err)
	}
	return nil
}

func (r *TransactionRepository) Update(ctx context.Context, tx *entity.Transaction) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE transactions SET status = $1, risk_score = $2, updated_at = $3
		WHERE id = $4`,
		string(tx.Status), tx.RiskScore, tx.UpdatedAt, tx.ID,
	)
	if err != nil {
		return fmt.Errorf("updating transaction: %w", err)
	}
	return nil
}

func (r *TransactionRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Transaction, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, account_id, amount, currency, merchant_id, location, status, risk_score, metadata, created_at, updated_at
		FROM transactions WHERE id = $1`, id)

	return scanTransaction(row)
}

func (r *TransactionRepository) FindByAccountID(ctx context.Context, accountID uuid.UUID) ([]*entity.Transaction, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, account_id, amount, currency, merchant_id, location, status, risk_score, metadata, created_at, updated_at
		FROM transactions WHERE account_id = $1 ORDER BY created_at DESC`, accountID)
	if err != nil {
		return nil, fmt.Errorf("querying transactions by account: %w", err)
	}
	defer rows.Close()

	var txs []*entity.Transaction
	for rows.Next() {
		tx, err := scanTransaction(rows)
		if err != nil {
			return nil, err
		}
		txs = append(txs, tx)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating transaction rows: %w", err)
	}
	return txs, nil
}

func (r *TransactionRepository) CountRecentByAccountID(ctx context.Context, accountID uuid.UUID, withinMinutes int) (int, error) {
	since := time.Now().UTC().Add(-time.Duration(withinMinutes) * time.Minute)
	var count int
	err := r.pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM transactions
		WHERE account_id = $1 AND created_at >= $2`, accountID, since).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("counting recent transactions: %w", err)
	}
	return count, nil
}

type scanner interface {
	Scan(dest ...any) error
}

func scanTransaction(s scanner) (*entity.Transaction, error) {
	var tx entity.Transaction
	var status string
	var metadata []byte

	err := s.Scan(
		&tx.ID, &tx.AccountID, &tx.Amount, &tx.Currency,
		&tx.MerchantID, &tx.Location, &status, &tx.RiskScore,
		&metadata, &tx.CreatedAt, &tx.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("scanning transaction row: %w", err)
	}

	tx.Status = entity.TransactionStatus(status)

	if len(metadata) > 0 {
		if err := json.Unmarshal(metadata, &tx.Metadata); err != nil {
			return nil, fmt.Errorf("unmarshaling metadata: %w", err)
		}
	}

	return &tx, nil
}
