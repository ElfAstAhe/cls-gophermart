package repository

import (
	"context"

	"github.com/ElfAstAhe/cls-gophermart/internal/bll/model"
)

type WithdrawRepository interface {
	Find(ctx context.Context, id string) (*model.Withdraw, error)
	ListByAccount(ctx context.Context, accountID string) ([]*model.Withdraw, error)
	GetWithdrawsByAccount(ctx context.Context, accountID string) (float64, error)
	Create(ctx context.Context, accountID string, withdraw *model.Withdraw) (*model.Withdraw, error)
}
