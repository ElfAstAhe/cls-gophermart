package db

import (
	"database/sql"
	"time"

	_cfg "github.com/ElfAstAhe/cls-gophermart.git/internal/app/config"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type postgresqlDB struct {
	DB     *sql.DB
	DBKind string
	Dsn    string
}

var db *postgresqlDB

func newPostgresqlDB(dsn string) (*postgresqlDB, error) {
	if db != nil {
		return db, nil
	}
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
		DBKind: _cfg.DBKindPostgres,
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
