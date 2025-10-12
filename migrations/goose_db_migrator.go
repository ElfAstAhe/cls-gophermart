package migrations

import (
	"context"
	"database/sql"

	_log "github.com/ElfAstAhe/cls-gophermart/internal/app/logger"
	_err "github.com/ElfAstAhe/cls-gophermart/pkg/error"
	"github.com/pressly/goose/v3"
)

// GooseDBMigrator is implementation of DBMigrator interface
type GooseDBMigrator struct {
	DB  *sql.DB
	ctx context.Context
	log _log.AppLogger
}

func NewGooseDBMigrator(ctx context.Context, db *sql.DB, logger _log.AppLogger) (*GooseDBMigrator, error) {
	return &GooseDBMigrator{
		DB:  db,
		ctx: ctx,
		log: logger,
	}, nil
}

// DBMigrator

func (g *GooseDBMigrator) Initialize() error {
	if err := goose.SetDialect("postgres"); err != nil {
		return _err.NewDBMigrationError("error select dialect", err)
	}
	goose.SetTableName("goose_version_history")
	goose.SetLogger(_log.NewGooseLogger(g.log))

	return nil
}

func (g *GooseDBMigrator) Up() error {
	if err := goose.UpContext(g.ctx, g.DB, ".", goose.WithAllowMissing()); err != nil {
		return _err.NewDBMigrationError("error migrate up", err)
	}
	return nil
}

// ==============
