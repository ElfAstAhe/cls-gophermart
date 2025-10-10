package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	_mod "github.com/ElfAstAhe/cls-gophermart/internal/bll/model"
	_repo "github.com/ElfAstAhe/cls-gophermart/internal/bll/repository"
	_utl "github.com/ElfAstAhe/cls-gophermart/internal/utils"
	_err "github.com/ElfAstAhe/cls-gophermart/pkg/error"
)

type AccountServiceImpl struct {
	accountRepo  _repo.AccountRepository
	withdrawRepo _repo.WithdrawRepository
	orderRepo    _repo.OrderRepository
}

func NewAccountServiceImpl(accountRepo _repo.AccountRepository, withdrawRepo _repo.WithdrawRepository, orderRepo _repo.OrderRepository) *AccountServiceImpl {
	return &AccountServiceImpl{
		accountRepo:  accountRepo,
		withdrawRepo: withdrawRepo,
		orderRepo:    orderRepo,
	}
}

func (a *AccountServiceImpl) GetUserBalance(ctx context.Context, userID string) (*_mod.AccountBalance, error) {
	accountID, err := a.validateAndGetAccount(ctx, userID)
	if err != nil {
		return nil, err
	}

	balance, err := a.accountRepo.GetBalance(ctx, accountID)
	if err != nil {
		return nil, err
	}

	return balance, nil
}

func (a *AccountServiceImpl) Withdraw(ctx context.Context, userID string, orderNumber string, withdrawAmount float64) error {
	if err := a.validateWithdraw(userID, orderNumber, withdrawAmount); err != nil {
		return err
	}

	accountID, err := a.validateAndGetAccount(ctx, userID)
	if err != nil {
		return err
	}

	balance, err := a.accountRepo.GetBalance(ctx, accountID)
	if err != nil {
		return err
	}
	if withdrawAmount > balance.Balance {
		return _err.NewBllInsufficientBalanceError(balance.Balance, withdrawAmount)
	}

	model := _mod.NewWithdraw(orderNumber, withdrawAmount, time.Now())
	if _, err := a.withdrawRepo.Create(ctx, accountID, model); err != nil {
		return err
	}

	return nil
}

func (a *AccountServiceImpl) validateWithdraw(userID string, orderNumber string, withdrawAmount float64) error {
	if strings.TrimSpace(userID) == "" {
		return _err.NewAuthUnauthorizedError("userID is null")
	}
	if err := _utl.ValidateOrderNumberByLuhn(orderNumber); err != nil {
		return err
	}
	if withdrawAmount <= 0.0 {
		return _err.NewAppInvalidArgumentError("WithdrawAmount", withdrawAmount)
	}

	return nil
}

func (a *AccountServiceImpl) ListAllWithdrawalsByUser(ctx context.Context, userID string) ([]*_mod.Withdraw, error) {
	accountID, err := a.validateAndGetAccount(ctx, userID)
	if err != nil {
		return nil, err
	}

	return a.withdrawRepo.ListByAccount(ctx, accountID)
}

func (a *AccountServiceImpl) CreateOrder(ctx context.Context, userID string, orderNumber string) error {
	accountID, err := a.validateAndGetAccount(ctx, userID)
	if err != nil {
		return err
	}

	if err := _utl.ValidateOrderNumberByLuhn(orderNumber); err != nil {
		return err
	}

	model := _mod.NewOrder(orderNumber, _mod.OrderStatusNew, 0.0, time.Now())
	if _, err := a.orderRepo.Create(ctx, accountID, model); err != nil {
		return err
	}

	// ToDo: add order into poll channel

	return nil
}

func (a *AccountServiceImpl) ListAllOrdersByUser(ctx context.Context, userID string) ([]*_mod.Order, error) {
	accountID, err := a.validateAndGetAccount(ctx, userID)
	if err != nil {
		return nil, err
	}

	return a.orderRepo.ListByAccount(ctx, accountID)
}

func (a *AccountServiceImpl) validateAndGetAccount(ctx context.Context, userID string) (string, error) {
	if strings.TrimSpace(userID) == "" {
		return "", _err.NewAuthUnauthorizedError("userID is null")
	}
	account, err := a.accountRepo.FindByUser(ctx, userID)
	if err != nil {
		return "", err
	}
	if account == nil {
		return "", _err.NewModelNotExistsError("account", fmt.Sprintf("userID=[%s]", userID))
	}

	return account.ID, nil
}
