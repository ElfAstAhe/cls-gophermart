package facade

import (
	"context"

	_dto "github.com/ElfAstAhe/cls-gophermart/internal/ep/dto/v1"
)

type UsersFacade interface {
	GetBalance(ctx context.Context) (*_dto.BalanceDto, error)
	ListOrders(ctx context.Context) ([]*_dto.OrderDto, error)
	ListWithdrawals(ctx context.Context) ([]*_dto.WithdrawDto, error)
	CreateOrder(ctx context.Context, orderNum []byte) error
	CreateWithdraw(ctx context.Context, request *_dto.WithdrawDto) error
}
