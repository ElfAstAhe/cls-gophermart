package facade

import (
	"context"
	"io"

	"github.com/ElfAstAhe/cls-gophermart/internal/ep/dto/v1"
)

type UsersFacade interface {
	GetBalance(ctx context.Context) (*v1.BalanceDto, error)
	ListOrders(ctx context.Context) ([]*v1.OrderDto, error)
	ListWithdrawals(ctx context.Context) ([]*v1.WithdrawDto, error)
	CreateOrder(ctx context.Context, orderNum io.Reader) error
	CreateWithdraw(ctx context.Context, withdraw io.Reader) error
}
