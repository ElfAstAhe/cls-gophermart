package mapper

import (
	"github.com/ElfAstAhe/cls-gophermart/internal/bll/model"
	"github.com/ElfAstAhe/cls-gophermart/internal/ep/dto/v1"
)

func ToBalanceDto(accBalance *model.AccountBalance) *v1.BalanceDto {
	if accBalance == nil {
		return nil
	}

	return v1.NewBalanceDto(accBalance.Balance, accBalance.WithdrawalsAmount)
}

func ToOrderDto(order *model.Order) *v1.OrderDto {
	if order == nil {
		return nil
	}

	return v1.NewOrderDto(order.Number, toOrderStatusDto(order.Status), order.AccrualAmount, order.UploadedAt)
}

func ToOrderDtoList(orders []*model.Order) []*v1.OrderDto {
	res := make([]*v1.OrderDto, 0)
	for _, order := range orders {
		res = append(res, ToOrderDto(order))
	}

	return res
}

func toOrderStatusDto(orderStatus string) v1.OrderStatus {
	switch orderStatus {
	case model.OrderStatusNew:
		return v1.OrderStatusNew
	case model.OrderStatusProcessing:
		return v1.OrderStatusProcessing
	case model.OrderStatusInvalid:
		return v1.OrderStatusInvalid
	case model.OrderStatusProcessed:
		return v1.OrderStatusProcessed
	}

	return ""
}

func ToWithdrawDto(withdraw *model.Withdraw) *v1.WithdrawDto {
	if withdraw == nil {
		return nil
	}

	return v1.NewWithdrawDto(withdraw.OrderNumber, withdraw.WithdrawAmount, withdraw.ProcessedAt)
}

func ToWithdrawDtoList(withdrawals []*model.Withdraw) []*v1.WithdrawDto {
	res := make([]*v1.WithdrawDto, 0)
	for _, withdraw := range withdrawals {
		res = append(res, ToWithdrawDto(withdraw))
	}

	return res
}
