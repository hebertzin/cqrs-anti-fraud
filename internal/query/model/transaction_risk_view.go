package model

import (
	"time"

	"github.com/google/uuid"
)

type RiskLevel string

const (
	RiskLevelLow    RiskLevel = "low"
	RiskLevelMedium RiskLevel = "medium"
	RiskLevelHigh   RiskLevel = "high"
)

type TransactionRiskView struct {
	ID           uuid.UUID `json:"id"`
	AccountID    uuid.UUID `json:"account_id"`
	Amount       float64   `json:"amount"`
	Currency     string    `json:"currency"`
	MerchantID   string    `json:"merchant_id"`
	Location     string    `json:"location"`
	Status       string    `json:"status"`
	RiskScore    float64   `json:"risk_score"`
	RiskLevel    RiskLevel `json:"risk_level"`
	FraudReasons []string  `json:"fraud_reasons,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func RiskLevelFromScore(score float64) RiskLevel {
	switch {
	case score >= 0.8:
		return RiskLevelHigh
	case score >= 0.5:
		return RiskLevelMedium
	default:
		return RiskLevelLow
	}
}
