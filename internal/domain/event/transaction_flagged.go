package event

import "github.com/google/uuid"

type TransactionFlagged struct {
	Base
	TransactionID uuid.UUID `json:"transaction_id"`
	AccountID     uuid.UUID `json:"account_id"`
	Reason        string    `json:"reason"`
	FlaggedBy     string    `json:"flagged_by"`
}

func NewTransactionFlagged(transactionID, accountID uuid.UUID, reason, flaggedBy string) TransactionFlagged {
	return TransactionFlagged{
		Base:          NewBase(TypeTransactionFlagged, transactionID),
		TransactionID: transactionID,
		AccountID:     accountID,
		Reason:        reason,
		FlaggedBy:     flaggedBy,
	}
}
