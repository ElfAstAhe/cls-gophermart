package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	_db "github.com/ElfAstAhe/cls-gophermart/internal/app/config/db"
	_mod "github.com/ElfAstAhe/cls-gophermart/internal/bll/model"
	_utl "github.com/ElfAstAhe/cls-gophermart/internal/utils"
	_err "github.com/ElfAstAhe/cls-gophermart/pkg/error"
	"github.com/google/uuid"
)

const (
	pgFindWithdrawSql          string = `select id, order_number, withdraw_amount, processed_at from withdrawals where id = $1;`
	pgListWithdrawByAccountSql string = `select id, order_number, withdraw_amount, processed_at from withdrawals where account_id = $1;`
	pgGetWithdrawsByAccountSql string = `select sum(withdraw_amount) from withdrawals where account_id = $1;`
	pgCreateWithdrawSql        string = `insert into withdrawals(id, account_id, order_number, withdraw_amount, processed_at) values ($1, $2, $3, $4, $5);`
)

type WithdrawPgRepository struct {
	db _db.DB
}

func NewWithdrawPgRepository(db _db.DB) *WithdrawPgRepository {
	return &WithdrawPgRepository{
		db: db,
	}
}

func (wr *WithdrawPgRepository) Find(ctx context.Context, id string) (*_mod.Withdraw, error) {
	if strings.TrimSpace(id) == "" {
		return nil, nil
	}

	model := _mod.Withdraw{}
	row := wr.db.GetDB().QueryRowContext(ctx, pgFindWithdrawSql, id)
	err := row.Scan(&model.ID, &model.OrderNumber, &model.WithdrawAmount, &model.ProcessedAt)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &model, nil
}

func (wr *WithdrawPgRepository) ListByAccount(ctx context.Context, accountID string) ([]*_mod.Withdraw, error) {
	res := make([]*_mod.Withdraw, 0)
	if strings.TrimSpace(accountID) == "" {
		return res, nil
	}

	rows, err := wr.db.GetDB().QueryContext(ctx, pgListWithdrawByAccountSql, accountID)
	if err != nil {
		return nil, err
	}
	defer _utl.CloseOnly(rows)

	for rows.Next() {
		var model _mod.Withdraw
		err := rows.Scan(&model.ID, &model.OrderNumber, &model.WithdrawAmount, &model.ProcessedAt)
		if err != nil && errors.Is(err, sql.ErrNoRows) {
			return res, nil
		} else if err != nil {
			return nil, err
		}

		res = append(res, &model)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return res, nil
}

func (wr *WithdrawPgRepository) GetWithdrawsByAccount(ctx context.Context, accountID string) (float64, error) {
	row := wr.db.GetDB().QueryRowContext(ctx, pgGetWithdrawsByAccountSql, accountID)
	var withdrawsSum float64
	err := row.Scan(&withdrawsSum)
	if err != nil {
		return 0, err
	}

	return withdrawsSum, nil
}

func (wr *WithdrawPgRepository) Create(ctx context.Context, accountID string, withdraw *_mod.Withdraw) (*_mod.Withdraw, error) {
	if err := wr.validateInstance(withdraw); err != nil {
		return nil, err
	}

	withdraw.ID = uuid.New().String()
	_, err := wr.db.GetDB().ExecContext(ctx, pgCreateWithdrawSql, withdraw.ID, accountID, withdraw.OrderNumber, withdraw.WithdrawAmount, withdraw.ProcessedAt)
	if err != nil {
		return nil, err
	}

	return withdraw, nil
}

func (wr *WithdrawPgRepository) validateInstance(instance *_mod.Withdraw) error {
	if instance == nil {
		return _err.NewModelValidationError("withdraw", "withdraw is null", nil)
	}
	if !(instance.WithdrawAmount > 0.0) {
		return _err.NewModelValidationError("withdraw", "withdraw amount must be greater zero", nil)
	}
	if strings.TrimSpace(instance.OrderNumber) == "" {
		return _err.NewModelValidationError("withdraw", "order number is null", nil)
	}

	return nil
}
