package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

const (
	pgCreateTableOrdersSQL string = `create table if not exists orders (
    id varchar(50) not null,
    account_id varchar(50) not null,
    doc_number varchar(100) not null,
    status varchar(50) not null,
    accrual_amount numeric(12,2) not null default 0.0,
    uploaded_at timestamptz not null default now(),
    constraint orders_pk primary key (id),
    constraint orders_uk unique (doc_number),
    constraint orders_fk_account foreign key (account_id) references accounts(id) on delete cascade,
    constraint orders_ch_status check (status in ('NEW', 'INVALID', 'PROCESSING', 'PROCESSED'))
);`
	pgDropTableOrdersSQL          string = `drop table if exists orders cascade;`
	pgCreateIndexOrdersAccountSQL string = `create index if not exists orders_idx_account on orders(account_id asc, uploaded_at desc);`
	pgDropIndexOrdersAccountSQL   string = `drop index if exists orders_idx_account cascade;`
	pgCreateIndexOrdersBalanceSQL string = `create index if not exists orders_idx_balance on orders(account_id, status) where status = 'PROCESSED';`
	pgDropIndexOrdersBalanceSQL   string = `drop index if exists orders_idx_balance cascade;`
)

func init() {
	goose.AddMigrationNoTxContext(up00003, down00003)
}

func up00003(ctx context.Context, db *sql.DB) error {
	if err := createTableOrders(ctx, db); err != nil {
		return err
	}
	if err := createIndexOrdersBalance(ctx, db); err != nil {
		return err
	}

	return createIndexOrdersAccount(ctx, db)
}

func down00003(ctx context.Context, db *sql.DB) error {
	if err := dropIndexOrdersAccount(ctx, db); err != nil {
		return err
	}
	if err := dropIndexOrdersBalance(ctx, db); err != nil {
		return err
	}

	return dropTableOrders(ctx, db)
}

func createTableOrders(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, pgCreateTableOrdersSQL)
	if err != nil {
		return err
	}

	return nil
}

func createIndexOrdersAccount(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, pgCreateIndexOrdersAccountSQL)
	if err != nil {
		return err
	}

	return nil
}

func createIndexOrdersBalance(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, pgCreateIndexOrdersBalanceSQL)
	if err != nil {
		return err
	}

	return nil
}

func dropTableOrders(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, pgDropTableOrdersSQL)
	if err != nil {
		return err
	}

	return nil
}

func dropIndexOrdersAccount(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, pgDropIndexOrdersAccountSQL)
	if err != nil {
		return err
	}

	return nil
}

func dropIndexOrdersBalance(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, pgDropIndexOrdersBalanceSQL)
	if err != nil {
		return err
	}

	return nil
}
