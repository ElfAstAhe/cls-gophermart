package model

import (
	"time"
)

type Order struct {
	Number        string    `db:"doc_number"`
	Status        string    `db:"status"`
	AccrualAmount float32   `db:"accrual_amount"`
	UploadedAt    time.Time `db:"uploaded_at"`
	//    AccountID     string    `db:"account_id"`
}
