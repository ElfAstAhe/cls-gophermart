package migrations

import (
	"context"
	"database/sql"

	"github.com/ElfAstAhe/cls-gophermart/internal/app/logger"
	errs "github.com/ElfAstAhe/cls-gophermart/pkg/error"
	"github.com/pressly/goose/v3"
)

// GooseDBMigrator is implementation of DBMigrator interface
type GooseDBMigrator struct {
	DB  *sql.DB
	ctx context.Context
	log logger.Logger
}

func NewGooseDBMigrator(ctx context.Context, db *sql.DB, logger logger.Logger) (*GooseDBMigrator, error) {
	return &GooseDBMigrator{
		DB:  db,
		ctx: ctx,
		log: logger,
	}, nil
}

// DBMigrator

func (g *GooseDBMigrator) Initialize() error {
	if err := goose.SetDialect("postgres"); err != nil {
		return errs.NewDBMigrationError("error select dialect", err)
	}
	goose.SetTableName("goose_version_history")
	goose.SetLogger(logger.NewGooseLogger(g.log))

	return nil
}

func (g *GooseDBMigrator) Up() error {
	if err := goose.UpContext(g.ctx, g.DB, ".", goose.WithAllowMissing()); err != nil {
		return errs.NewDBMigrationError("error migrate up", err)
	}
	return nil
}

// ==============
