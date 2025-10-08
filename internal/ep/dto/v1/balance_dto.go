package v1

type BalanceDto struct {
	BalanceAmount float64 `json:"current"`
	WithdrawSum   float64 `json:"withdrawn"`
}

func NewBalanceDto(balance float64, withdrawSum float64) *BalanceDto {
	return &BalanceDto{
		BalanceAmount: balance,
		WithdrawSum:   withdrawSum,
	}
}
