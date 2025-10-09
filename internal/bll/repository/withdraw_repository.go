package repository

import (
	"context"

	_mod "github.com/ElfAstAhe/cls-gophermart/internal/bll/model"
)

type WithdrawRepository interface {
	Find(ctx context.Context, id string) (*_mod.Withdraw, error)
	ListByAccount(ctx context.Context, accountID string) ([]*_mod.Withdraw, error)
	GetWithdrawsByAccount(ctx context.Context, accountID string) (float64, error)
	Create(ctx context.Context, accountID string, withdraw *_mod.Withdraw) (*_mod.Withdraw, error)
}
