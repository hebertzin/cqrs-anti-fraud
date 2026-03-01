package model

import "github.com/google/uuid"

const AnalyzeTransactionCommand = "AnalyzeTransaction"

type AnalyzeTransaction struct {
	AccountID  uuid.UUID         `json:"account_id"`
	Amount     float64           `json:"amount"`
	Currency   string            `json:"currency"`
	MerchantID string            `json:"merchant_id"`
	Location   string            `json:"location"`
	Metadata   map[string]string `json:"metadata,omitempty"`
}

type AnalyzeTransactionResult struct {
	TransactionID uuid.UUID `json:"transaction_id"`
	Status        string    `json:"status"`
	RiskScore     float64   `json:"risk_score"`
	RiskLevel     string    `json:"risk_level"`
}
