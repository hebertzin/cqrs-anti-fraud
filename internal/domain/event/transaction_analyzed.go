package event

import (
	"github.com/google/uuid"
	"github.com/hebertzin/cqrs/internal/domain/entity"
)

type TransactionAnalyzed struct {
	Base
	TransactionID uuid.UUID               `json:"transaction_id"`
	AccountID     uuid.UUID               `json:"account_id"`
	Amount        float64                 `json:"amount"`
	Currency      string                  `json:"currency"`
	MerchantID    string                  `json:"merchant_id"`
	Location      string                  `json:"location"`
	RiskScore     float64                 `json:"risk_score"`
	Status        entity.TransactionStatus `json:"status"`
	FraudReasons  []string                `json:"fraud_reasons,omitempty"`
}

func NewTransactionAnalyzed(t *entity.Transaction, reasons []string) TransactionAnalyzed {
	return TransactionAnalyzed{
		Base:          NewBase(TypeTransactionAnalyzed, t.ID),
		TransactionID: t.ID,
		AccountID:     t.AccountID,
		Amount:        t.Amount,
		Currency:      t.Currency,
		MerchantID:    t.MerchantID,
		Location:      t.Location,
		RiskScore:     t.RiskScore,
		Status:        t.Status,
		FraudReasons:  reasons,
	}
}
