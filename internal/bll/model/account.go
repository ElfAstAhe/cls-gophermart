package model

type Account struct {
	ID          string `db:"id"`
	Person      string `db:"person"`
	Orders      []*Order
	Withdrawals []*Withdraw
}

func NewAccount(id string, person string) *Account {
	return &Account{
		ID:          id,
		Person:      person,
		Orders:      make([]*Order, 0),
		Withdrawals: make([]*Withdraw, 0),
	}
}
