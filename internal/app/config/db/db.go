package db

import (
	"database/sql"
	"io"
)

type DB interface {
	GetDB() *sql.DB
	GetDBKind() string
	GetDsn() string
}

func NewDB(kind string, dsn string) (DB, error) {
	return newPostgresqlDB(dsn)
}

func CloseDB(db DB) error {
	if closer, ok := db.(io.Closer); ok {
		return closer.Close()
	}

	return nil
}
