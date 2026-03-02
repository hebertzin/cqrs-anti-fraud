package model

import "github.com/google/uuid"

const FlagTransactionCommand = "FlagTransaction"

type FlagTransaction struct {
	TransactionID uuid.UUID `json:"transaction_id"`
	Reason        string    `json:"reason"`
	FlaggedBy     string    `json:"flagged_by"`
}
