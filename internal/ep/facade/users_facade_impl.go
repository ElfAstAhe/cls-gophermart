package facade

import (
	"context"
	"encoding/json"
	"io"

	"github.com/ElfAstAhe/cls-gophermart/internal/app/logger"
	"github.com/ElfAstAhe/cls-gophermart/internal/bll/service"
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
	userID := ctx.Value("user_id")
	if userID == nil {
		return nil, errs.NewAuthUnauthorizedError("user_id not found in context")
	}
	// get balance model
	model, err := ufi.accountService.GetUserBalance(ctx, userID.(string))
	if err != nil {
		return nil, err
	}
	// transform
	dto := mapper.ToBalanceDto(model)

	return dto, nil
}

func (ufi *UsersFacadeImpl) ListOrders(ctx context.Context) ([]*v1.OrderDto, error) {
	// get user id
	userID := ctx.Value("user_id")
	if userID == nil {
		return nil, errs.NewAuthUnauthorizedError("user_id not found in context")
	}
	// get orders
	modelList, err := ufi.accountService.ListAllOrdersByUser(ctx, userID.(string))
	if err != nil {
		return nil, err
	}
	// transform
	dtoList := mapper.ToOrderDtoList(modelList)

	return dtoList, nil
}

func (ufi *UsersFacadeImpl) ListWithdrawals(ctx context.Context) ([]*v1.WithdrawDto, error) {
	userID := ctx.Value("user_id")
	if userID == nil {
		return nil, errs.NewAuthUnauthorizedError("user_id not found in context")
	}
	modelList, err := ufi.accountService.ListAllWithdrawalsByUser(ctx, userID.(string))
	if err != nil {
		return nil, err
	}
	dtoList := mapper.ToWithdrawDtoList(modelList)

	return dtoList, nil
}

func (ufi *UsersFacadeImpl) CreateOrder(ctx context.Context, orderNum io.Reader) error {
	userID := ctx.Value("user_id")
	if userID == nil {
		return errs.NewAuthUnauthorizedError("user_id not found in context")
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
		return errs.NewAuthUnauthorizedError("user_id not found in context")
	}
	dec := json.NewDecoder(withdraw)
	dto := &v1.WithdrawDto{}
	err := dec.Decode(dto)
	if err != nil {
		return err
	}

	return ufi.accountService.Withdraw(ctx, userID.(string), dto.Order, dto.WithdrawAmount)
}
