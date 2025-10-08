package model

type Account struct {
	ID            string  `db:"id"`
	AccrualAmount float32 `db:"accrual_amount"`
	Orders        []*Order
	Withdrawals   []*Withdraw
}
