package service

import (
	"context"

	_mod "github.com/ElfAstAhe/cls-gophermart/internal/bll/model"
)

type AccountService interface {
	GetUserBalance(ctx context.Context, userID string) (*_mod.AccountBalance, error)
	Withdraw(ctx context.Context, userID string, orderNumber string, accrualAmount float64) error
	ListAllWithdrawalsByUser(ctx context.Context, userID string) ([]*_mod.Withdraw, error)
	CreateOrder(ctx context.Context, userID string, orderNumber string) error
	ListAllOrdersByUser(ctx context.Context, userID string) ([]*_mod.Order, error)
}
