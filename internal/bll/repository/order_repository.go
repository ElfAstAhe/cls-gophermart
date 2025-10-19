package repository

import (
	"context"

	"github.com/ElfAstAhe/cls-gophermart/internal/bll/model"
)

type OrderRepository interface {
	Find(ctx context.Context, id string) (*model.Order, error)
	FindByNumber(ctx context.Context, number string) (*model.Order, error)
	ListByAccount(ctx context.Context, accountId string) ([]*model.Order, error)
	ListUnfinished(ctx context.Context, statuses ...string) ([]*model.Order, error)
	Create(ctx context.Context, accountId string, order *model.Order) (*model.Order, error)
	Change(ctx context.Context, order *model.Order) (*model.Order, error)
}
