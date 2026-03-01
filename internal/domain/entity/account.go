package entity

import (
	"time"

	"github.com/google/uuid"
)

type AccountStatus string

const (
	AccountStatusActive  AccountStatus = "active"
	AccountStatusBlocked AccountStatus = "blocked"
	AccountStatusFlagged AccountStatus = "flagged"
)

type Account struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Status    AccountStatus
	RiskLevel float64
	BlockedAt *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewAccount(userID uuid.UUID) *Account {
	now := time.Now().UTC()
	return &Account{
		ID:        uuid.New(),
		UserID:    userID,
		Status:    AccountStatusActive,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (a *Account) Block() {
	now := time.Now().UTC()
	a.Status = AccountStatusBlocked
	a.BlockedAt = &now
	a.UpdatedAt = now
}

func (a *Account) Flag() {
	a.Status = AccountStatusFlagged
	a.UpdatedAt = time.Now().UTC()
}

func (a *Account) Activate() {
	a.Status = AccountStatusActive
	a.BlockedAt = nil
	a.UpdatedAt = time.Now().UTC()
}

func (a *Account) IsBlocked() bool {
	return a.Status == AccountStatusBlocked
}

func (a *Account) UpdateRiskLevel(level float64) {
	a.RiskLevel = level
	a.UpdatedAt = time.Now().UTC()
}
