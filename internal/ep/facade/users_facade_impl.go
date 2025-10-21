package facade

import (
	"context"
	"encoding/json"
	"io"

	"github.com/ElfAstAhe/cls-gophermart/internal/app/logger"
	"github.com/ElfAstAhe/cls-gophermart/internal/bll/service"
	"github.com/ElfAstAhe/cls-gophermart/internal/bll/service/auth"
	"github.com/ElfAstAhe/cls-gophermart/internal/ep/dto/v1"
	"github.com/ElfAstAhe/cls-gophermart/internal/ep/mapper"
	errs "github.com/ElfAstAhe/cls-gophermart/pkg/error"
)

type UsersFacadeImpl struct {
	accountService service.AccountService
	log            logger.Logger
}

func NewUsersFacadeImpl(accountSvc service.AccountService, logger logger.Logger) *UsersFacadeImpl {
	return &UsersFacadeImpl{
		accountService: accountSvc,
		log:            logger.GetLogger("UsersFacadeImpl"),
	}
}

func (ufi *UsersFacadeImpl) GetBalance(ctx context.Context) (*v1.BalanceDto, error) {
	// get user id
	userInfo, err := auth.UserInfoFromContext(ctx)
	if err != nil {
		return nil, errs.NewAuthUnauthorizedError("userInfo not found in context", err)
	}
	// get balance model
	model, err := ufi.accountService.GetUserBalance(ctx, userInfo.UserID)
	if err != nil {
		return nil, err
	}
	// transform
	dto := mapper.ToBalanceDto(model)

	return dto, nil
}

func (ufi *UsersFacadeImpl) ListOrders(ctx context.Context) ([]*v1.OrderDto, error) {
	// get user id
	userInfo, err := auth.UserInfoFromContext(ctx)
	if err != nil {
		return nil, errs.NewAuthUnauthorizedError("userInfo not found in context", err)
	}
	// get orders
	modelList, err := ufi.accountService.ListAllOrdersByUser(ctx, userInfo.UserID)
	if err != nil {
		return nil, err
	}
	// transform
	dtoList := mapper.ToOrderDtoList(modelList)

	return dtoList, nil
}

func (ufi *UsersFacadeImpl) ListWithdrawals(ctx context.Context) ([]*v1.WithdrawDto, error) {
	userInfo, err := auth.UserInfoFromContext(ctx)
	if err != nil {
		return nil, errs.NewAuthUnauthorizedError("userInfo not found in context", err)
	}
	modelList, err := ufi.accountService.ListAllWithdrawalsByUser(ctx, userInfo.UserID)
	if err != nil {
		return nil, err
	}
	dtoList := mapper.ToWithdrawDtoList(modelList)

	return dtoList, nil
}

func (ufi *UsersFacadeImpl) CreateOrder(ctx context.Context, orderNum io.Reader) error {
	userInfo, err := auth.UserInfoFromContext(ctx)
	if err != nil {
		return errs.NewAuthUnauthorizedError("userInfo not found in context", err)
	}
	buf, err := io.ReadAll(orderNum)
	if err != nil {
		return err
	}

	return ufi.accountService.CreateOrder(ctx, userInfo.UserID, string(buf))
}

func (ufi *UsersFacadeImpl) CreateWithdraw(ctx context.Context, withdraw io.Reader) error {
	userInfo, err := auth.UserInfoFromContext(ctx)
	if err != nil {
		return errs.NewAuthUnauthorizedError("userInfo not found in context", err)
	}
	dec := json.NewDecoder(withdraw)
	dto := &v1.WithdrawDto{}
	err = dec.Decode(dto)
	if err != nil {
		return err
	}

	return ufi.accountService.Withdraw(ctx, userInfo.UserID, dto.Order, dto.WithdrawAmount)
}
