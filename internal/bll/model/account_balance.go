package model

type AccountBalance struct {
	ID                string  `db:"account_id"`
	Balance           float64 `db:"balance"`
	AccrualsAmount    float64 `db:"accruals_amount"`
	WithdrawalsAmount float64 `db:"withdrawals_amount"`
}

func NewAccountBalance(id string, balance float64, accrualsAmount float64, withdrawalsAmount float64) *AccountBalance {
	return &AccountBalance{
		ID:                id,
		Balance:           balance,
		AccrualsAmount:    accrualsAmount,
		WithdrawalsAmount: withdrawalsAmount,
	}
}
