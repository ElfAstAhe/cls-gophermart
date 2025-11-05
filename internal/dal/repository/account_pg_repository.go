package repository

import (
	"context"
	"errors"
	"strings"

	"github.com/ElfAstAhe/cls-gophermart/internal/app/config/db"
	"github.com/ElfAstAhe/cls-gophermart/internal/bll/model"
	"github.com/ElfAstAhe/cls-gophermart/internal/bll/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

const (
	pgFindAccountSQL       string = `select id, person from accounts where id = $1`
	pgFindAccountByUserSQL string = `select id, person from accounts where user_id = $1`
	pgCreateAccountSQL     string = `insert into accounts (id, user_id, person) values ($1, $2, $3)`
	pgChangeAccountSQL     string = `update accounts set person = $2 where id = $1`
	pgGetAccountBalanceSQL string = `select account_id, balance, accruals_amount, withdrawals_amount from v_accounts_balance where account_id = $1`
)

type AccountPgRepository struct {
	db           db.DB
	withdrawRepo repository.WithdrawRepository
	orderRepo    repository.OrderRepository
}

func NewAccountPgRepository(db db.DB, withdrawRepository repository.WithdrawRepository, orderRepository repository.OrderRepository) *AccountPgRepository {
	return &AccountPgRepository{
		db:           db,
		withdrawRepo: withdrawRepository,
		orderRepo:    orderRepository,
	}
}

func (a *AccountPgRepository) Find(ctx context.Context, id string) (*model.Account, error) {
	if strings.TrimSpace(id) == "" {
		return nil, nil
	}

	return a.findSingle(ctx, pgFindAccountSQL, id)
}

func (a *AccountPgRepository) FindByUser(ctx context.Context, userID string) (*model.Account, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, nil
	}

	return a.findSingle(ctx, pgFindAccountByUserSQL, userID)
}

func (a *AccountPgRepository) findSingle(ctx context.Context, query string, param any) (*model.Account, error) {
	row := a.db.GetDB().QueryRowContext(ctx, query, param)

	account := &model.Account{}
	err := row.Scan(&account.ID, &account.Person)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	withdrawals, err := a.withdrawRepo.ListByAccount(ctx, account.ID)
	if err != nil {
		return nil, err
	}
	account.Withdrawals = withdrawals
	orders, err := a.orderRepo.ListByAccount(ctx, account.ID)
	if err != nil {
		return nil, err
	}
	account.Orders = orders

	return account, nil
}

func (a *AccountPgRepository) Create(ctx context.Context, userID string, account *model.Account) (*model.Account, error) {
	account.ID = uuid.New().String()

	_, err := a.db.GetDB().ExecContext(ctx, pgCreateAccountSQL, account.ID, userID, account.Person)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (a *AccountPgRepository) Change(ctx context.Context, userID string, account *model.Account) (*model.Account, error) {
	_, err := a.db.GetDB().ExecContext(ctx, pgChangeAccountSQL, account.ID, account.Person)
	if err != nil {
		return nil, err
	}

	return account, nil
}

// GetBalance ToDo: maybe need to move into separate repository
// GetBalance get account current full balance info
func (a *AccountPgRepository) GetBalance(ctx context.Context, id string) (*model.AccountBalance, error) {
	row := a.db.GetDB().QueryRowContext(ctx, pgGetAccountBalanceSQL, id)
	if row.Err() != nil {
		return nil, row.Err()
	}
	balance := &model.AccountBalance{}
	err := row.Scan(&balance.ID, &balance.Balance, &balance.AccrualsAmount, &balance.WithdrawalsAmount)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return balance, nil
}
