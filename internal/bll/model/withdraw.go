package model

import (
	"time"
)

type Withdraw struct {
	ID string `db:"id"`
	//    AccountID      string    `db:"account_id"`
	OrderNumber    string    `db:"order_number"`
	WithdrawAmount float32   `db:"withdraw_amount"`
	ProcessedAt    time.Time `db:"processed_at"`
}
