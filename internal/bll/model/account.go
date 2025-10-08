package model

type Account struct {
	ID            string  `db:"id"`
	AccrualAmount float64 `db:"accrual_amount"`
	Orders        []*Order
	Withdrawals   []*Withdraw
}

func NewAccount() *Account {
	return &Account{}
}
