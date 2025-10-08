package service

import (
	"context"

	_mod "github.com/ElfAstAhe/cls-gophermart/internal/bll/model"
)

type OrdersService interface {
	Create(ctx context.Context, userID string, orderNumber string) error
	ListAllByUser(ctx context.Context, userID string) ([]*_mod.Order, error)
}
