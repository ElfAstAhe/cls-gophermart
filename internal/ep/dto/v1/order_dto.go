package v1

import "time"

type OrderStatus string

const (
	OrderStatusNew        OrderStatus = "NEW"
	OrderStatusProcessing OrderStatus = "PROCESSING"
	OrderStatusInvalid    OrderStatus = "INVALID"
	OrderStatusProcessed  OrderStatus = "PROCESSED"
)

type OrderDto struct {
	DocNumber     string      `json:"number"`
	Status        OrderStatus `json:"status"`
	AccrualAmount float64     `json:"accrual,omitempty"`
	UploadedAt    time.Time   `json:"uploaded_at"`
}

func NewOrderDto(docNumber string, status OrderStatus, accrualAmount float64, uploadedAt time.Time) *OrderDto {
	return &OrderDto{
		DocNumber:     docNumber,
		Status:        status,
		AccrualAmount: accrualAmount,
		UploadedAt:    uploadedAt,
	}
}
