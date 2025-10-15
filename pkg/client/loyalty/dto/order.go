package dto

const (
	LSOrderStatusRegistered string = "REGISTERED"
	LSOrderStatusInvalid    string = "INVALID"
	LSOrderStatusProcessing string = "PROCESSING"
	LSOrderStatusProcessed  string = "PROCESSED"
)

type LSOrderDto struct {
	Number        string  `json:"order"`
	Status        string  `json:"status"`
	AccrualAmount float64 `json:"accrual,omitempty"`
}
