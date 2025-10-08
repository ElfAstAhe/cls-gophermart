package model

import (
	"time"
)

type Withdraw struct {
	ID             string    `db:"id"`
	OrderNumber    string    `db:"order_number"`
	WithdrawAmount float64   `db:"withdraw_amount"`
	ProcessedAt    time.Time `db:"processed_at"`
}

func NewWithdraw(orderNumber string, withdrawAmount float64, processedAt time.Time) *Withdraw {
	return &Withdraw{
		OrderNumber:    orderNumber,
		WithdrawAmount: withdrawAmount,
		ProcessedAt:    processedAt,
	}
}
