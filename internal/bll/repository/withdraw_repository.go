package repository

import (
	"context"

	_mod "github.com/ElfAstAhe/cls-gophermart/internal/bll/model"
)

type WithdrawRepository interface {
	Find(ctx context.Context, id string) (*_mod.Withdraw, error)
	ListByUser(ctx context.Context, userId string) ([]*_mod.Withdraw, error)
	WithdrawsByUser(ctx context.Context, userId string) (float64, error)
}
