package migrations

import (
	"context"
	"database/sql"

	_log "github.com/ElfAstAhe/cls-gophermart/internal/app/logger"
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
		log: logger.GetLogger("migration"),
	}, nil
}

// DBMigrator

func (g *GooseDBMigrator) Initialize() error {
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	goose.SetTableName("goose_version_history")
	goose.SetLogger(_log.NewGooseLogger(g.log))

	return nil
}

func (g *GooseDBMigrator) Up() error {
	return goose.UpContext(g.ctx, g.DB, ".", goose.WithAllowMissing())
}

// ==============
