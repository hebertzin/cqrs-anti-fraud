package model

import (
	"time"

	"github.com/google/uuid"
)

type AccountStatusView struct {
	ID                uuid.UUID  `json:"id"`
	UserID            uuid.UUID  `json:"user_id"`
	Status            string     `json:"status"`
	RiskLevel         float64    `json:"risk_level"`
	TotalTransactions int        `json:"total_transactions"`
	FlaggedCount      int        `json:"flagged_count"`
	DeclinedCount     int        `json:"declined_count"`
	BlockedAt         *time.Time `json:"blocked_at,omitempty"`
	LastActivityAt    *time.Time `json:"last_activity_at,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
}
