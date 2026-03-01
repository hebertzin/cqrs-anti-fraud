package event

import "github.com/google/uuid"

type AccountBlocked struct {
	Base
	AccountID uuid.UUID `json:"account_id"`
	Reason    string    `json:"reason"`
	BlockedBy string    `json:"blocked_by"`
}

func NewAccountBlocked(accountID uuid.UUID, reason, blockedBy string) AccountBlocked {
	return AccountBlocked{
		Base:      NewBase(TypeAccountBlocked, accountID),
		AccountID: accountID,
		Reason:    reason,
		BlockedBy: blockedBy,
	}
}
