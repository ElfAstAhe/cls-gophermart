package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

const (
	pgCreateTableWithdrawalsSQL string = `create table if not exists withdrawals (
    id varchar(50) not null,
    account_id varchar(50) not null,
    order_number varchar(100) not null,
    withdraw_amount numeric(12, 2) not null default 0.0,
    processed_at timestamptz not null default now(),
    constraint withdrawals_pk primary key (id),
    constraint withdrawals_fk foreign key (account_id) references accounts(id) on delete cascade,
    constraint withdrawals_ch_amount check (withdraw_amount >= 0.0)
);`
	pgDropTableWithdrawalsSQL   string = `drop table if exists withdrawals cascade;`
	pgCreateIndexWithdrawalsSQL string = `create index if not exists withdrawals_idx on withdrawals(account_id asc, processed_at desc);`
	pgDropIndexWithdrawalsSQL   string = `drop index if exists withdrawals_idx cascade;`
)

func init() {
	goose.AddMigrationNoTxContext(up00004, down00004)
}

func up00004(ctx context.Context, db *sql.DB) error {
	if err := createTableWithdrawals(ctx, db); err != nil {
		return err
	}

	return createIndexWithdrawals(ctx, db)
}

func down00004(ctx context.Context, db *sql.DB) error {
	if err := dropIndexWithdrawals(ctx, db); err != nil {
		return err
	}

	return dropTableWithdrawals(ctx, db)
}

func createTableWithdrawals(ctx context.Context, db *sql.DB) error {
	_, err := db.Exec(pgCreateTableWithdrawalsSQL)
	if err != nil {
		return err
	}

	return nil
}

func createIndexWithdrawals(ctx context.Context, db *sql.DB) error {
	_, err := db.Exec(pgCreateIndexWithdrawalsSQL)
	if err != nil {
		return err
	}

	return nil
}

func dropTableWithdrawals(ctx context.Context, db *sql.DB) error {
	_, err := db.Exec(pgDropTableWithdrawalsSQL)
	if err != nil {
		return err
	}

	return nil
}

func dropIndexWithdrawals(ctx context.Context, db *sql.DB) error {
	_, err := db.Exec(pgDropIndexWithdrawalsSQL)
	if err != nil {
		return err
	}

	return nil
}
