package repository

import (
    "context"
    "database/sql"
    "errors"
    "fmt"
    "slices"
    "strings"

    "github.com/ElfAstAhe/cls-gophermart/internal/app/config/db"
    "github.com/ElfAstAhe/cls-gophermart/internal/bll/model"
    "github.com/ElfAstAhe/cls-gophermart/internal/utils"
    errs "github.com/ElfAstAhe/cls-gophermart/pkg/error"
    "github.com/google/uuid"
)

const (
    pgFindOrderSql            string = `select id, account_id, doc_number, status, accrual_amount, uploaded_at from orders where id = $1`
    pgFindOrderByNumberSql    string = `select id, account_id, doc_number, status, accrual_amount, uploaded_at from orders where doc_number = $1`
    pgListOrdersByAccountSql  string = `select id, account_id, doc_number, status, accrual_amount, uploaded_at from orders where account_id = $1`
    pgCreateOrderSql          string = `insert into orders(id, account_id, doc_number, status, accrual_amount, uploaded_at) values($1, $2, $3, $4, $5, $6)`
    pgChangeOrderSql          string = `update orders set status = $2, accrual_amount = $3, uploaded_at = $4 where id = $1`
    pgListOrdersUnfinishedSql string = `select id, account_id, doc_number, status, accrual_amount, uploaded_at from orders where status = any($1)`
)

type OrderPgRepository struct {
    db db.DB
}

func NewOrderPgRepository(db db.DB) *OrderPgRepository {
    return &OrderPgRepository{
        db: db,
    }
}

func (o *OrderPgRepository) listAcceptableOrderStatuses() []string {
    return []string{
        model.OrderStatusNew,
        model.OrderStatusProcessing,
        model.OrderStatusInvalid,
        model.OrderStatusProcessed,
    }
}

func (o *OrderPgRepository) listUnacceptableOrderStatuses() []string {
    return []string{
        model.OrderStatusInvalid,
        model.OrderStatusProcessed,
    }
}

func (o *OrderPgRepository) Find(ctx context.Context, id string) (*model.Order, error) {
    return o.findSingle(ctx, pgFindOrderSql, id)
}

func (o *OrderPgRepository) FindByNumber(ctx context.Context, number string) (*model.Order, error) {
    return o.findSingle(ctx, pgFindOrderByNumberSql, number)
}

func (o *OrderPgRepository) findSingle(ctx context.Context, query string, params ...any) (*model.Order, error) {
    row := o.db.GetDB().QueryRowContext(ctx, query, params...)

    order := &model.Order{}
    err := row.Scan(&order.ID, &order.AccountID, &order.Number, &order.Status, &order.AccrualAmount, &order.UploadedAt)
    if err != nil && errors.Is(err, sql.ErrNoRows) {
        return nil, nil
    } else if err != nil {
        return nil, err
    }

    return order, nil
}

func (o *OrderPgRepository) ListByAccount(ctx context.Context, accountId string) ([]*model.Order, error) {
    if accountId == "" {
        return make([]*model.Order, 0), nil
    }

    return o.list(ctx, pgListOrdersByAccountSql, accountId)
}

func (o *OrderPgRepository) ListUnfinished(ctx context.Context, statuses ...string) ([]*model.Order, error) {
    if len(statuses) == 0 {
        return make([]*model.Order, 0), nil
    }

    return o.list(ctx, pgListOrdersUnfinishedSql, statuses)
}

func (o *OrderPgRepository) list(ctx context.Context, query string, params ...any) ([]*model.Order, error) {
    res := make([]*model.Order, 0)
    rows, err := o.db.GetDB().QueryContext(ctx, query, params...)
    if err != nil {
        return nil, err
    }
    defer utils.CloseOnly(rows)

    for rows.Next() {
        order := &model.Order{}

        err := rows.Scan(&order.ID, &order.AccountID, &order.Number, &order.Status, &order.AccrualAmount, &order.UploadedAt)
        if err != nil && errors.Is(err, sql.ErrNoRows) {
            return res, nil
        } else if err != nil {
            return nil, err
        }

        res = append(res, order)
    }
    if rows.Err() != nil {
        return nil, rows.Err()
    }

    return res, nil
}

func (o *OrderPgRepository) Create(ctx context.Context, accountId string, order *model.Order) (*model.Order, error) {
    if err := o.validateInstance(order); err != nil {
        return nil, err
    }

    if orderExists, err := o.validateCreateIntegrity(ctx, order); err != nil {
        return orderExists, err
    }

    order.ID = uuid.New().String()

    _, err := o.db.GetDB().ExecContext(ctx, pgCreateOrderSql, order.ID, accountId, order.Number, order.Status, order.AccrualAmount, order.UploadedAt)
    if err != nil {
        return nil, err
    }

    return order, nil
}

func (o *OrderPgRepository) Change(ctx context.Context, order *model.Order) (*model.Order, error) {
    if err := o.validateInstance(order); err != nil {
        return nil, err
    }
    if err := o.validateChangeIntegrity(ctx, order); err != nil {
        return nil, err
    }

    _, err := o.db.GetDB().ExecContext(ctx, pgChangeOrderSql, order.ID, order.Status, order.AccrualAmount, order.UploadedAt)
    if err != nil {
        return nil, err
    }

    return order, nil
}

func (o *OrderPgRepository) validateInstance(order *model.Order) error {
    // nil instance
    if order == nil {
        return errs.NewModelValidationError("order", "order is null", nil)
    }

    // order number
    if strings.TrimSpace(order.Number) == "" {
        return errs.NewModelValidationError("order", "order number empty", nil)
    }

    // order status
    if err := o.validateInstanceStatus(order); err != nil {
        return errs.NewModelValidationError("order", "order status invalid", err)
    }

    // amount
    if !(order.AccrualAmount >= 0.0) {
        return errs.NewModelValidationError("order", "order accrual_amount must be greater than zero", nil)
    }

    return nil
}

func (o *OrderPgRepository) validateInstanceStatus(order *model.Order) error {
    if !slices.Contains(o.listAcceptableOrderStatuses(), order.Status) {
        return fmt.Errorf("order status [%v] invalid", order.Status)
    }

    return nil
}

func (o *OrderPgRepository) validateCreateIntegrity(ctx context.Context, order *model.Order) (*model.Order, error) {
    orderExists, err := o.FindByNumber(ctx, order.Number)
    if err != nil {
        return nil, err
    }
    if orderExists != nil {
        return orderExists, errs.NewModelAlreadyExistsError("order", order.Number)
    }

    return nil, nil
}

func (o *OrderPgRepository) validateChangeIntegrity(ctx context.Context, order *model.Order) error {
    // model existence
    founded, err := o.Find(ctx, order.ID)
    if err != nil {
        return err
    }
    if founded == nil {
        return errs.NewModelNotExistsError("order", order.ID)
    }

    return nil
}
