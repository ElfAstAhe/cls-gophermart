package loyalty

import (
	_dto "github.com/ElfAstAhe/cls-gophermart/pkg/client/loyalty/dto"
)

type LSClient interface {
	GetOrder(orderNumber string) (*_dto.LSOrderDto, error)
}
