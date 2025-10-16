package facade

import (
	"context"

	_dto "github.com/ElfAstAhe/cls-gophermart/internal/ep/dto/v1"
)

type UsersFacade interface {
	GetBalance(ctx context.Context) (*_dto.BalanceDto, error)
	GetOrders(ctx context.Context) ([]*_dto.OrderDto, error)
}
