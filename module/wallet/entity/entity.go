package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

type User struct {
	ID             uint32
	Name           string
	Balance        decimal.Decimal
	PendingBalance decimal.Decimal
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type Disbursement struct {
	ID        uint32
	UserID    uint32
	Amount    decimal.Decimal
	Status    DisbursementStatusEnum
	CreatedAt time.Time
	UpdatedAt time.Time
}

type DisburseRequest struct {
	UserID uint32          `json:"user_id"`
	Amount decimal.Decimal `json:"amount"`
}

type DisburseResponse struct {
	DisbursementID         uint32
	DisbursementStatusEnum DisbursementStatusEnum
	Message                string
}
