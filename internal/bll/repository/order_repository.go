package repository

import (
	"context"

	_mod "github.com/ElfAstAhe/cls-gophermart/internal/bll/model"
)

type OrderRepository interface {
	Find(ctx context.Context, id string) (*_mod.Order, error)
	FindByNumber(ctx context.Context, number string) (*_mod.Order, error)
}
