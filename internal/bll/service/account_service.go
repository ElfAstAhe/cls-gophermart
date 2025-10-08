package service

import (
	"context"

	_mod "github.com/ElfAstAhe/cls-gophermart/internal/bll/model"
)

type AccountService interface {
	GetUserBalance(ctx context.Context, userID string) (float64, float64, error)
	Withdraw(ctx context.Context, userID string, orderNumber string, accrualAmount float64) error
	ListAllByUser(ctx context.Context, userID string) ([]*_mod.Withdraw, error)
}
