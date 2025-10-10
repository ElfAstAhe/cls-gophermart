package repository

import (
	"context"
	"errors"
	"strings"

	_db "github.com/ElfAstAhe/cls-gophermart/internal/app/config/db"
	_mod "github.com/ElfAstAhe/cls-gophermart/internal/bll/model"
	_rep "github.com/ElfAstAhe/cls-gophermart/internal/bll/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

const (
	pgFindAccountSql       string = `select id, person from accounts where id = $1`
	pgFindAccountByUserSql string = `select id, person from accounts where user_id = $1`
	pgCreateAccountSql     string = `insert into accounts (id, user_id, person) values ($1, $2, $3)`
	pgChangeAccountSql     string = `update accounts set person = $2 where id = $1`
	pgGetAccountBalanceSql string = `select account_id, balance, accruals_amount, withdrawals_amount from v_accounts_balance where account_id = $1`
)

type AccountPgRepository struct {
	db           _db.DB
	withdrawRepo _rep.WithdrawRepository
	orderRepo    _rep.OrderRepository
}

func NewAccountPgRepository(db _db.DB, withdrawRepository _rep.WithdrawRepository, orderRepository _rep.OrderRepository) *AccountPgRepository {
	return &AccountPgRepository{
		db:           db,
		withdrawRepo: withdrawRepository,
		orderRepo:    orderRepository,
	}
}

func (a *AccountPgRepository) Find(ctx context.Context, id string) (*_mod.Account, error) {
	if strings.TrimSpace(id) == "" {
		return nil, nil
	}

	return a.findSingle(ctx, pgFindAccountSql, id)
}

func (a *AccountPgRepository) FindByUser(ctx context.Context, userID string) (*_mod.Account, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, nil
	}

	return a.findSingle(ctx, pgFindAccountByUserSql, userID)
}

func (a *AccountPgRepository) findSingle(ctx context.Context, query string, param any) (*_mod.Account, error) {
	row := a.db.GetDB().QueryRowContext(ctx, query, param)

	model := &_mod.Account{}
	err := row.Scan(&model.ID, &model.Person)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	withdrawals, err := a.withdrawRepo.ListByAccount(ctx, model.ID)
	if err != nil {
		return nil, err
	}
	model.Withdrawals = withdrawals
	orders, err := a.orderRepo.ListByAccount(ctx, model.ID)
	if err != nil {
		return nil, err
	}
	model.Orders = orders

	return model, nil
}

func (a *AccountPgRepository) Create(ctx context.Context, userID string, account *_mod.Account) (*_mod.Account, error) {
	account.ID = uuid.New().String()

	_, err := a.db.GetDB().ExecContext(ctx, pgCreateAccountSql, account.ID, userID, account.Person)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (a *AccountPgRepository) Change(ctx context.Context, userID string, account *_mod.Account) (*_mod.Account, error) {
	_, err := a.db.GetDB().ExecContext(ctx, pgChangeAccountSql, account.ID, account.Person)
	if err != nil {
		return nil, err
	}

	return account, nil
}

// GetBalance ToDo: maybe need to move into separate repository
// GetBalance get account current full balance info
func (a *AccountPgRepository) GetBalance(ctx context.Context, id string) (*_mod.AccountBalance, error) {
	row := a.db.GetDB().QueryRowContext(ctx, pgGetAccountBalanceSql, id)
	if row.Err() != nil {
		return nil, row.Err()
	}
	model := &_mod.AccountBalance{}
	err := row.Scan(&model.ID, &model.Balance, &model.AccrualsAmount, &model.WithdrawalsAmount)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return model, nil
}
