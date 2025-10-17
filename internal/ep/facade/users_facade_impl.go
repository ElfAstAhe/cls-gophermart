package facade

import (
	"context"
	"encoding/json"
	"io"

	_log "github.com/ElfAstAhe/cls-gophermart/internal/app/logger"
	_svc "github.com/ElfAstAhe/cls-gophermart/internal/bll/service"
	_dto "github.com/ElfAstAhe/cls-gophermart/internal/ep/dto/v1"
	_map "github.com/ElfAstAhe/cls-gophermart/internal/ep/mapper"
	_err "github.com/ElfAstAhe/cls-gophermart/pkg/error"
)

type UsersFacadeImpl struct {
	accountService _svc.AccountService
	log            _log.AppLogger
}

func NewUsersFacadeImpl(accountSvc _svc.AccountService, logger _log.AppLogger) *UsersFacadeImpl {
	return &UsersFacadeImpl{
		accountService: accountSvc,
		log:            logger.GetLogger("UsersFacadeImpl"),
	}
}

func (ufi *UsersFacadeImpl) GetBalance(ctx context.Context) (*_dto.BalanceDto, error) {
	// get user id
	userID := ctx.Value("user_id")
	if userID == nil {
		return nil, _err.NewAuthUnauthorizedError("user_id not found in context")
	}
	// get balance model
	model, err := ufi.accountService.GetUserBalance(ctx, userID.(string))
	if err != nil {
		return nil, err
	}
	// transform
	dto := _map.ToBalanceDto(model)

	return dto, nil
}

func (ufi *UsersFacadeImpl) ListOrders(ctx context.Context) ([]*_dto.OrderDto, error) {
	// get user id
	userID := ctx.Value("user_id")
	if userID == nil {
		return nil, _err.NewAuthUnauthorizedError("user_id not found in context")
	}
	// get orders
	modelList, err := ufi.accountService.ListAllOrdersByUser(ctx, userID.(string))
	if err != nil {
		return nil, err
	}
	// transform
	dtoList := _map.ToOrderDtoList(modelList)

	return dtoList, nil
}

func (ufi *UsersFacadeImpl) ListWithdrawals(ctx context.Context) ([]*_dto.WithdrawDto, error) {
	userID := ctx.Value("user_id")
	if userID == nil {
		return nil, _err.NewAuthUnauthorizedError("user_id not found in context")
	}
	modelList, err := ufi.accountService.ListAllWithdrawalsByUser(ctx, userID.(string))
	if err != nil {
		return nil, err
	}
	dtoList := _map.ToWithdrawDtoList(modelList)

	return dtoList, nil
}

func (ufi *UsersFacadeImpl) CreateOrder(ctx context.Context, orderNum io.Reader) error {
	userID := ctx.Value("user_id")
	if userID == nil {
		return _err.NewAuthUnauthorizedError("user_id not found in context")
	}
	buf, err := io.ReadAll(orderNum)
	if err != nil {
		return err
	}

	return ufi.accountService.CreateOrder(ctx, userID.(string), string(buf))
}

func (ufi *UsersFacadeImpl) CreateWithdraw(ctx context.Context, withdraw io.Reader) error {
	userID := ctx.Value("user_id")
	if userID == nil {
		return _err.NewAuthUnauthorizedError("user_id not found in context")
	}
	dec := json.NewDecoder(withdraw)
	dto := &_dto.WithdrawDto{}
	err := dec.Decode(dto)
	if err != nil {
		return err
	}

	return ufi.accountService.Withdraw(ctx, userID.(string), dto.Order, dto.WithdrawAmount)
}
