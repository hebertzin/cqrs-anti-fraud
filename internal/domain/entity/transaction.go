package entity

import (
	"time"

	"github.com/google/uuid"
)

type TransactionStatus string

const (
	TransactionStatusPending  TransactionStatus = "pending"
	TransactionStatusApproved TransactionStatus = "approved"
	TransactionStatusDeclined TransactionStatus = "declined"
	TransactionStatusFlagged  TransactionStatus = "flagged"
)

type Transaction struct {
	ID         uuid.UUID
	AccountID  uuid.UUID
	Amount     float64
	Currency   string
	MerchantID string
	Location   string
	Status     TransactionStatus
	RiskScore  float64
	Metadata   map[string]string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func NewTransaction(accountID uuid.UUID, amount float64, currency, merchantID, location string) *Transaction {
	now := time.Now().UTC()
	return &Transaction{
		ID:         uuid.New(),
		AccountID:  accountID,
		Amount:     amount,
		Currency:   currency,
		MerchantID: merchantID,
		Location:   location,
		Status:     TransactionStatusPending,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

func (t *Transaction) Approve() {
	t.Status = TransactionStatusApproved
	t.UpdatedAt = time.Now().UTC()
}

func (t *Transaction) Decline() {
	t.Status = TransactionStatusDeclined
	t.UpdatedAt = time.Now().UTC()
}

func (t *Transaction) Flag() {
	t.Status = TransactionStatusFlagged
	t.UpdatedAt = time.Now().UTC()
}

func (t *Transaction) SetRiskScore(score float64) {
	t.RiskScore = score
	t.UpdatedAt = time.Now().UTC()
}

func (t *Transaction) IsHighRisk() bool {
	return t.RiskScore >= 0.8
}

func (t *Transaction) IsMediumRisk() bool {
	return t.RiskScore >= 0.5 && t.RiskScore < 0.8
}
