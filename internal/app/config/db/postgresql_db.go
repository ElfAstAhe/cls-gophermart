package db

import (
	"database/sql"
	"time"

	"github.com/ElfAstAhe/cls-gophermart/internal/app/config"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type postgresqlDB struct {
	DB     *sql.DB
	DBKind string
	Dsn    string
}

func newPostgresqlDB(dsn string) (*postgresqlDB, error) {
	pg, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	pg.SetMaxOpenConns(20)
	pg.SetMaxIdleConns(5)
	pg.SetConnMaxIdleTime(60 * time.Second)

	err = pg.Ping()
	if err != nil {
		return nil, err
	}

	return &postgresqlDB{
		DB:     pg,
		DBKind: config.DBKindPostgres,
		Dsn:    dsn,
	}, nil
}

// Closer

func (pDB *postgresqlDB) Close() error {
	return pDB.DB.Close()
}

// =============

// DB

func (pDB *postgresqlDB) GetDB() *sql.DB {
	return pDB.DB
}

func (pDB *postgresqlDB) GetDBKind() string {
	return pDB.DBKind
}

func (pDB *postgresqlDB) GetDsn() string {
	return pDB.Dsn
}

// =============
