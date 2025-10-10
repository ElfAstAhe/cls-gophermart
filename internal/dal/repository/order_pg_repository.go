package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"

	_db "github.com/ElfAstAhe/cls-gophermart/internal/app/config/db"
	_mod "github.com/ElfAstAhe/cls-gophermart/internal/bll/model"
	_utl "github.com/ElfAstAhe/cls-gophermart/internal/utils"
	_err "github.com/ElfAstAhe/cls-gophermart/pkg/error"
	"github.com/google/uuid"
)

const (
	pgFindOrderSql            string = `select id, doc_number, status, accrual_amount, uploaded_at from orders where id = $1`
	pgFindOrderByNumberSql    string = `select id, doc_number, status, accrual_amount, uploaded_at from orders where doc_number = $1`
	pgListOrdersByAccountSql  string = `select id, doc_number, status, accrual_amount, uploaded_at from orders where account_id = $1`
	pgCreateOrderSql          string = `insert into orders(id, account_id, doc_number, status, accrual_amount, uploaded_at) values($1, $2, $3, $4, $5, $6)`
	pgChangeOrderSql          string = `update orders set status = $2, accrual_amount = $3, uploaded_at = $4 where id = $1`
	pgListOrdersUnfinishedSql string = `select id, doc_number, status, accrual_amount, uploaded_at from orders where status = any($1)`
)

type OrderPgRepository struct {
	db _db.DB
}

var changeUnacceptableOrderStatuses []string = []string{
	_mod.OrderStatusInvalid,
	_mod.OrderStatusProcessed,
}
var acceptableOrderStatuses []string = []string{
	_mod.OrderStatusNew,
	_mod.OrderStatusProcessing,
	_mod.OrderStatusInvalid,
	_mod.OrderStatusProcessed,
}

func NewOrderPgRepository(db _db.DB) *OrderPgRepository {
	return &OrderPgRepository{
		db: db,
	}
}

func (o *OrderPgRepository) Find(ctx context.Context, id string) (*_mod.Order, error) {
	return o.findSingle(ctx, pgFindOrderSql, id)
}

func (o *OrderPgRepository) FindByNumber(ctx context.Context, number string) (*_mod.Order, error) {
	return o.findSingle(ctx, pgFindOrderByNumberSql, number)
}

func (o *OrderPgRepository) findSingle(ctx context.Context, query string, param any) (*_mod.Order, error) {
	row := o.db.GetDB().QueryRowContext(ctx, query, param)

	model := &_mod.Order{}
	err := row.Scan(&model.ID, &model.Number, &model.Status, &model.AccrualAmount, &model.UploadedAt)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return model, nil
}

func (o *OrderPgRepository) ListByAccount(ctx context.Context, accountId string) ([]*_mod.Order, error) {
	if accountId == "" {
		return make([]*_mod.Order, 0), nil
	}

	return o.list(ctx, pgListOrdersByAccountSql, accountId)
}

func (o *OrderPgRepository) ListUnfinished(ctx context.Context, statuses ...string) ([]*_mod.Order, error) {
	if len(statuses) == 0 {
		return make([]*_mod.Order, 0), nil
	}

	return o.list(ctx, pgListOrdersUnfinishedSql, statuses)
}

func (o *OrderPgRepository) list(ctx context.Context, query string, params ...any) ([]*_mod.Order, error) {
	res := make([]*_mod.Order, 0)
	rows, err := o.db.GetDB().QueryContext(ctx, query, params...)
	if err != nil {
		return nil, err
	}
	defer _utl.CloseOnly(rows)

	for rows.Next() {
		model := &_mod.Order{}

		err := rows.Scan(&model.ID, &model.Number, &model.Status, &model.AccrualAmount, &model.UploadedAt)
		if err != nil && errors.Is(err, sql.ErrNoRows) {
			return res, nil
		} else if err != nil {
			return nil, err
		}

		res = append(res, model)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return res, nil
}

func (o *OrderPgRepository) Create(ctx context.Context, accountId string, order *_mod.Order) (*_mod.Order, error) {
	if err := o.validateInstance(order); err != nil {
		return nil, err
	}

	if err := o.validateCreateBL(ctx, order); err != nil {
		return nil, err
	}

	order.ID = uuid.New().String()

	_, err := o.db.GetDB().ExecContext(ctx, pgCreateOrderSql, order.ID, accountId, order.Number, order.Status, order.AccrualAmount, order.UploadedAt)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (o *OrderPgRepository) Change(ctx context.Context, order *_mod.Order) (*_mod.Order, error) {
	if err := o.validateInstance(order); err != nil {
		return nil, err
	}
	if err := o.validateChangeBL(ctx, order); err != nil {
		return nil, err
	}

	_, err := o.db.GetDB().ExecContext(ctx, pgChangeOrderSql, order.ID, order.Status, order.AccrualAmount, order.UploadedAt)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (o *OrderPgRepository) validateInstance(order *_mod.Order) error {
	// nil instance
	if order == nil {
		return _err.NewModelValidationError("order", "order is null", nil)
	}

	// order number (luna algorithm)
	if err := _utl.ValidateOrderNumberByLuhn(order.Number); err != nil {
		return _err.NewModelValidationError("order", "order number invalid", err)
	}

	// order status
	if err := o.validateInstanceStatus(order); err != nil {
		return _err.NewModelValidationError("order", "order status invalid", err)
	}

	// amount
	if !(order.AccrualAmount >= 0.0) {
		return _err.NewModelValidationError("order", "order accrual_amount must be greater than zero", nil)
	}

	return nil
}

func (o *OrderPgRepository) validateInstanceStatus(order *_mod.Order) error {
	if !slices.Contains(acceptableOrderStatuses, order.Status) {
		return fmt.Errorf("order status [%v] invalid", order.Status)
	}

	return nil
}

func (o *OrderPgRepository) validateCreateBL(ctx context.Context, order *_mod.Order) error {
	model, err := o.FindByNumber(ctx, order.Number)
	if err != nil {
		return err
	}
	if model != nil {
		return _err.NewModelAlreadyExistsError("order", order.Number)
	}

	return nil
}

func (o *OrderPgRepository) validateChangeBL(ctx context.Context, order *_mod.Order) error {
	// model existence
	model, err := o.Find(ctx, order.ID)
	if err != nil {
		return err
	}
	if model == nil {
		return _err.NewModelNotExistsError("order", order.ID)
	}
	// acceptable status
	if slices.Contains(changeUnacceptableOrderStatuses, order.Status) {
		return _err.NewModelValidationError("order", "order status is invalid", nil)
	}
	// amount
	if !(order.AccrualAmount >= 0.0) {
		return _err.NewModelValidationError("order", "order accrual_amount must be equal or greater than zero", nil)
	}

	return nil
}
