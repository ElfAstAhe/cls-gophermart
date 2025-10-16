package mapper

import (
	_mod "github.com/ElfAstAhe/cls-gophermart/internal/bll/model"
	_dto "github.com/ElfAstAhe/cls-gophermart/internal/ep/dto/v1"
)

func ToBalanceDto(model *_mod.AccountBalance) *_dto.BalanceDto {
	if model == nil {
		return nil
	}

	return _dto.NewBalanceDto(model.Balance, model.WithdrawalsAmount)
}

func ToOrderDto(model *_mod.Order) *_dto.OrderDto {
	if model == nil {
		return nil
	}

	return _dto.NewOrderDto(model.Number, toOrderStatusDto(model.Status), model.AccrualAmount, model.UploadedAt)
}

func ToOrderDtoList(models []*_mod.Order) []*_dto.OrderDto {
	res := make([]*_dto.OrderDto, len(models))
	for i, model := range models {
		res[i] = ToOrderDto(model)
	}

	return res
}

func toOrderStatusDto(modelStatus string) _dto.OrderStatus {
	switch modelStatus {
	case _mod.OrderStatusNew:
		return _dto.OrderStatusNew
	case _mod.OrderStatusProcessing:
		return _dto.OrderStatusProcessing
	case _mod.OrderStatusInvalid:
		return _dto.OrderStatusInvalid
	case _mod.OrderStatusProcessed:
		return _dto.OrderStatusProcessed
	}

	return ""
}
