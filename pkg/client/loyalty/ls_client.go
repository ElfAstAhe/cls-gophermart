package loyalty

import (
	"context"

	_dto "github.com/ElfAstAhe/cls-gophermart/pkg/client/loyalty/dto"
)

type LSClient interface {
	GetOrder(ctx context.Context, orderNumber string) (*_dto.LSOrderDto, error)
}
