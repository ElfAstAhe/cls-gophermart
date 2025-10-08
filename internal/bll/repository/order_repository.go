package repository

import (
	"context"

	_mod "github.com/ElfAstAhe/cls-gophermart/internal/bll/model"
)

type OrderRepository interface {
	Find(ctx context.Context, id string) (*_mod.Order, error)
	FindByNumber(ctx context.Context, number string) (*_mod.Order, error)
	ListByAccount(ctx context.Context, accountId string) ([]*_mod.Order, error)
	Create(ctx context.Context, accountId string, order *_mod.Order) (*_mod.Order, error)
	Change(ctx context.Context, accountId string, order *_mod.Order) (*_mod.Order, error)
}
