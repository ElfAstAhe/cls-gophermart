package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

const createTableOrdersSql string = `create table if not exists orders (
    id varchar(50) not null,
    account_id varchar(50) not null,
    doc_number varchar(100) not null,
    status varchar(50) not null,
    accrual_amount numeric(12,2) not null default 0.0,
    uploaded_at timestamptz not null default now(),
    constraint orders_pk primary key (id),
    constraint orders_uk unique (doc_number),
    constraint orders_fk_account foreign key (account_id) references accounts(id) on delete cascade,
    constraint orders_ch_status check (status in ('REGISTERED', 'INVALID', 'PROCESSING', 'PROCESSED'))
);`
const dropTableOrdersSql string = `drop table if exists orders cascade;`
const createIndexOrdersSql string = `create index if not exists orders_idx_status on orders(status asc);`
const dropIndexOrdersSql string = `drop index if exists orders_idx_status cascade;`

func init() {
	goose.AddMigrationNoTxContext(up00003, down00003)
}

func up00003(ctx context.Context, db *sql.DB) error {
	if err := createTableOrders(ctx, db); err != nil {
		return err
	}

	return createIndexOrders(ctx, db)
}

func down00003(ctx context.Context, db *sql.DB) error {
	if err := dropIndexOrders(ctx, db); err != nil {
		return err
	}

	return dropTableOrders(ctx, db)
}

func createTableOrders(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, createTableOrdersSql)
	if err != nil {
		return err
	}

	return nil
}

func createIndexOrders(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, createIndexOrdersSql)
	if err != nil {
		return err
	}

	return nil
}

func dropTableOrders(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, dropTableOrdersSql)
	if err != nil {
		return err
	}

	return nil
}

func dropIndexOrders(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, dropIndexOrdersSql)
	if err != nil {
		return err
	}

	return nil
}
