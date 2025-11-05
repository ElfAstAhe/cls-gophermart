package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

const (
	pgCreateViewOrdersBalanceSql string = `create or replace view v_accounts_balance as
select
    a.id as account_id,
    coalesce(ao.accruals_amount, 0.0) - coalesce(aw.withdrawals_amount, 0.0) as balance,
    coalesce(ao.accruals_amount, 0.0) as accruals_amount,
    coalesce(aw.withdrawals_amount, 0.0) as withdrawals_amount
from
    accounts a
    left outer join (
        select
            account_id,
            sum(accrual_amount) as accruals_amount
        from
            orders
        where
            status = 'PROCESSED'
        group by
            account_id
    ) ao
        on
            ao.account_id = a.id
    left outer join (
        select
            account_id,
            sum(withdraw_amount) as withdrawals_amount
        from
            withdrawals
        group by
            account_id
    ) aw
        on
            aw.account_id = a.id
;`
	pgDropViewOrdersBalanceSql string = `drop view if exists v_orders_balance;`
)

func init() {
	goose.AddMigrationNoTxContext(up00005, down00005)
}

func up00005(ctx context.Context, db *sql.DB) error {
	return createViewOrdersBalance(ctx, db)
}

func down00005(ctx context.Context, db *sql.DB) error {
	return dropViewOrdersBalance(ctx, db)
}

func createViewOrdersBalance(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, pgCreateViewOrdersBalanceSql)
	if err != nil {
		return err
	}

	return nil
}

func dropViewOrdersBalance(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, pgDropViewOrdersBalanceSql)
	if err != nil {
		return err
	}

	return nil
}
