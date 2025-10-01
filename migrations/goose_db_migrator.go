package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

// GooseDBMigrator is implementation of DBMigrator interface
type GooseDBMigrator struct {
	DB  *sql.DB
	ctx context.Context
	log *zap.SugaredLogger
}

func NewGooseDBMigrator(ctx context.Context, db *sql.DB, logger *zap.Logger) (*GooseDBMigrator, error) {
	return &GooseDBMigrator{
		DB:  db,
		ctx: ctx,
		log: logger.Sugar(),
	}, nil
}

// DBMigrator

func (g *GooseDBMigrator) Initialize() error {
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	goose.SetTableName("goose_version_history")
	goose.SetLogger(newGooseLogger(g.log))

	return nil
}

func (g *GooseDBMigrator) Up() error {
	return goose.UpContext(g.ctx, g.DB, ".", goose.WithAllowMissing())
}

// ==============
