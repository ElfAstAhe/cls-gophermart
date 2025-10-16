package v1

import "time"

type WithdrawDto struct {
	Order          string    `json:"order"`
	WithdrawAmount float64   `json:"sum"`
	ProcessedAt    time.Time `json:"processed_at"`
}

func NewWithdrawDto(order string, accrualAmount float64, processedAt time.Time) *WithdrawDto {
	return &WithdrawDto{
		Order:          order,
		WithdrawAmount: accrualAmount,
		ProcessedAt:    processedAt,
	}
}
