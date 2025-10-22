package service

import (
    "context"
    "errors"
    "fmt"
    "strings"
    "time"

    "github.com/ElfAstAhe/cls-gophermart/internal/bll/model"
    "github.com/ElfAstAhe/cls-gophermart/internal/bll/repository"
    "github.com/ElfAstAhe/cls-gophermart/internal/utils"
    errs "github.com/ElfAstAhe/cls-gophermart/pkg/error"
)

type AccountServiceImpl struct {
    accountRepo      repository.AccountRepository
    withdrawRepo     repository.WithdrawRepository
    orderRepo        repository.OrderRepository
    orderPollService OrdersPollingService
}

func NewAccountServiceImpl(orderPollService OrdersPollingService, accountRepo repository.AccountRepository, withdrawRepo repository.WithdrawRepository, orderRepo repository.OrderRepository) *AccountServiceImpl {
    return &AccountServiceImpl{
        accountRepo:      accountRepo,
        withdrawRepo:     withdrawRepo,
        orderRepo:        orderRepo,
        orderPollService: orderPollService,
    }
}

func (a *AccountServiceImpl) GetUserBalance(ctx context.Context, userID string) (*model.AccountBalance, error) {
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
        return errs.NewBllInsufficientBalanceError(balance.Balance, withdrawAmount)
    }

    withdraw := model.NewWithdraw(orderNumber, withdrawAmount, time.Now())
    if _, err := a.withdrawRepo.Create(ctx, accountID, withdraw); err != nil {
        return err
    }

    return nil
}

func (a *AccountServiceImpl) validateWithdraw(userID string, orderNumber string, withdrawAmount float64) error {
    if strings.TrimSpace(userID) == "" {
        return errs.NewAuthUnauthorizedError("userID is null", nil)
    }
    if err := utils.ValidateOrderNumberByLuhn(orderNumber); err != nil {
        return err
    }
    if withdrawAmount <= 0.0 {
        return errs.NewAppInvalidArgumentError("WithdrawAmount", withdrawAmount)
    }

    return nil
}

func (a *AccountServiceImpl) ListAllWithdrawalsByUser(ctx context.Context, userID string) ([]*model.Withdraw, error) {
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

    if err := utils.ValidateOrderNumberByLuhn(orderNumber); err != nil {
        return err
    }

    order, err := a.orderRepo.Create(ctx, accountID, model.NewOrder(orderNumber, model.OrderStatusNew, 0.0, time.Now()))
    if err != nil {
        if errors.As(err, &errs.ModelAlreadyExistsErr) {
            if order.AccountID == accountID {
                return errs.NewBllOrderAlreadyExistsCurrentError(order.Number)
            }

            return errs.NewBllOrderAlreadyExistsAnotherError(order.Number)
        }

        return err
    }

    // add id for order poll
    a.orderPollService.Add(order.ID)

    return nil
}

func (a *AccountServiceImpl) ListAllOrdersByUser(ctx context.Context, userID string) ([]*model.Order, error) {
    accountID, err := a.validateAndGetAccount(ctx, userID)
    if err != nil {
        return nil, err
    }

    return a.orderRepo.ListByAccount(ctx, accountID)
}

func (a *AccountServiceImpl) validateAndGetAccount(ctx context.Context, userID string) (string, error) {
    if strings.TrimSpace(userID) == "" {
        return "", errs.NewAuthUnauthorizedError("userID is null", nil)
    }
    account, err := a.accountRepo.FindByUser(ctx, userID)
    if err != nil {
        return "", err
    }
    if account == nil {
        return "", errs.NewModelNotExistsError("account", fmt.Sprintf("userID=[%s]", userID))
    }

    return account.ID, nil
}
