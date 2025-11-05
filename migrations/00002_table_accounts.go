package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

const (
	pgCreateTableAccountsSQL string = `create table if not exists accounts(
    id varchar(50) not null,
    user_id varchar(50) not null,
    person varchar(256) null,
    constraint accounts_pk primary key(id),
    constraint accounts_uk unique (user_id),
    constraint accounts_fk_user foreign key (user_id) references users(id) on delete cascade
);`
	pgDropTableAccountsSQL string = `drop table if exists accounts cascade;`
)

func init() {
	goose.AddMigrationNoTxContext(up00002, down00002)
}

func up00002(ctx context.Context, db *sql.DB) error {
	return createTableAccounts(ctx, db)
}

func down00002(ctx context.Context, db *sql.DB) error {
	return dropTableAccounts(ctx, db)
}

func createTableAccounts(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, pgCreateTableAccountsSQL)
	if err != nil {
		return err
	}

	return nil
}

func dropTableAccounts(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, pgDropTableAccountsSQL)
	if err != nil {
		return err
	}

	return nil
}
