package model

import "github.com/google/uuid"

const BlockAccountCommand = "BlockAccount"

type BlockAccount struct {
	AccountID uuid.UUID `json:"account_id"`
	Reason    string    `json:"reason"`
	BlockedBy string    `json:"blocked_by"`
}
