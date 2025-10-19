package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/ElfAstAhe/cls-gophermart/internal/app/config/db"
	"github.com/ElfAstAhe/cls-gophermart/internal/bll/model"
	"github.com/ElfAstAhe/cls-gophermart/internal/utils"
	errs "github.com/ElfAstAhe/cls-gophermart/pkg/error"
	"github.com/google/uuid"
)

const (
	pgFindWithdrawSql          string = `select id, order_number, withdraw_amount, processed_at from withdrawals where id = $1;`
	pgListWithdrawByAccountSql string = `select id, order_number, withdraw_amount, processed_at from withdrawals where account_id = $1 order by processed_at desc;`
	pgGetWithdrawsByAccountSql string = `select sum(withdraw_amount) from withdrawals where account_id = $1;`
	pgCreateWithdrawSql        string = `insert into withdrawals(id, account_id, order_number, withdraw_amount, processed_at) values ($1, $2, $3, $4, $5);`
)

type WithdrawPgRepository struct {
	db db.DB
}

func NewWithdrawPgRepository(db db.DB) *WithdrawPgRepository {
	return &WithdrawPgRepository{
		db: db,
	}
}

func (wr *WithdrawPgRepository) Find(ctx context.Context, id string) (*model.Withdraw, error) {
	if strings.TrimSpace(id) == "" {
		return nil, nil
	}

	withdraw := model.Withdraw{}
	row := wr.db.GetDB().QueryRowContext(ctx, pgFindWithdrawSql, id)
	err := row.Scan(&withdraw.ID, &withdraw.OrderNumber, &withdraw.WithdrawAmount, &withdraw.ProcessedAt)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &withdraw, nil
}

func (wr *WithdrawPgRepository) ListByAccount(ctx context.Context, accountID string) ([]*model.Withdraw, error) {
	res := make([]*model.Withdraw, 0)
	if strings.TrimSpace(accountID) == "" {
		return res, nil
	}

	rows, err := wr.db.GetDB().QueryContext(ctx, pgListWithdrawByAccountSql, accountID)
	if err != nil {
		return nil, err
	}
	defer utils.CloseOnly(rows)

	for rows.Next() {
		var withdraw model.Withdraw
		err := rows.Scan(&withdraw.ID, &withdraw.OrderNumber, &withdraw.WithdrawAmount, &withdraw.ProcessedAt)
		if err != nil && errors.Is(err, sql.ErrNoRows) {
			return res, nil
		} else if err != nil {
			return nil, err
		}

		res = append(res, &withdraw)
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

func (wr *WithdrawPgRepository) Create(ctx context.Context, accountID string, withdraw *model.Withdraw) (*model.Withdraw, error) {
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

func (wr *WithdrawPgRepository) validateInstance(instance *model.Withdraw) error {
	if instance == nil {
		return errs.NewModelValidationError("withdraw", "withdraw is null", nil)
	}
	if !(instance.WithdrawAmount > 0.0) {
		return errs.NewModelValidationError("withdraw", "withdraw amount must be greater zero", nil)
	}
	if err := utils.ValidateOrderNumberByLuhn(instance.OrderNumber); err != nil {
		return err
	}

	return nil
}
