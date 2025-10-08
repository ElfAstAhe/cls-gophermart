package model

import (
	"time"
)

const (
	OrderStatusNew        string = "NEW"
	OrderStatusProcessing string = "PROCESSING"
	OrderStatusInvalid    string = "INVALID"
	OrderStatusProcessed  string = "PROCESSED"
)

type Order struct {
	ID            string    `db:"id"`
	Number        string    `db:"doc_number"`
	Status        string    `db:"status"`
	AccrualAmount float64   `db:"accrual_amount"`
	UploadedAt    time.Time `db:"uploaded_at"`
}

func NewOrder(number string, status string, accrualAmount float64, uploadedAt time.Time) *Order {
	return &Order{
		Number:        number,
		Status:        status,
		AccrualAmount: accrualAmount,
		UploadedAt:    uploadedAt,
	}
}
