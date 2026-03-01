package model

import (
	"time"

	"github.com/google/uuid"
)

type FraudAlertStatus string

const (
	FraudAlertStatusOpen     FraudAlertStatus = "open"
	FraudAlertStatusReviewed FraudAlertStatus = "reviewed"
	FraudAlertStatusDismissed FraudAlertStatus = "dismissed"
)

type FraudAlertView struct {
	ID            uuid.UUID        `json:"id"`
	TransactionID uuid.UUID        `json:"transaction_id"`
	AccountID     uuid.UUID        `json:"account_id"`
	Amount        float64          `json:"amount"`
	Currency      string           `json:"currency"`
	RiskScore     float64          `json:"risk_score"`
	RiskLevel     RiskLevel        `json:"risk_level"`
	Reasons       []string         `json:"reasons"`
	Status        FraudAlertStatus `json:"status"`
	CreatedAt     time.Time        `json:"created_at"`
}

type FraudAlertListResponse struct {
	Alerts []*FraudAlertView `json:"alerts"`
	Total  int               `json:"total"`
	Page   int               `json:"page"`
	Limit  int               `json:"limit"`
}
