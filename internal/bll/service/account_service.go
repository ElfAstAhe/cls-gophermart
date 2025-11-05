package service

import (
	"context"

	"github.com/ElfAstAhe/cls-gophermart/internal/bll/model"
)

type AccountService interface {
	GetUserBalance(ctx context.Context, userID string) (*model.AccountBalance, error)
	Withdraw(ctx context.Context, userID string, orderNumber string, accrualAmount float64) error
	ListAllWithdrawalsByUser(ctx context.Context, userID string) ([]*model.Withdraw, error)
	CreateOrder(ctx context.Context, userID string, orderNumber string) error
	ListAllOrdersByUser(ctx context.Context, userID string) ([]*model.Order, error)
}
