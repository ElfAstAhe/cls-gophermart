package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

const createTableUsersSql string = `create table if not exists users(
    id varchar(50) not null,
    username varchar(100) not null,
    password varchar(1024) not null,
    disabled bool not null default false,
    constraint users_pk primary key (id),
    constraint user_uk_username unique(username)
);`

const dropTableUsersSql string = `drop table if exists users cascade;`

func init() {
	goose.AddMigrationNoTxContext(up00001, down00001)
}

func up00001(ctx context.Context, db *sql.DB) error {
	//	return createTableUsers(ctx, db)
	return nil
}

func down00001(ctx context.Context, db *sql.DB) error {
	//	return dropTableUsers(ctx, db)
	return nil
}

func createTableUsers(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, createTableUsersSql)
	if err != nil {
		return err
	}

	return nil
}

func dropTableUsers(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, dropTableUsersSql)
	if err != nil {
		return err
	}

	return nil
}
