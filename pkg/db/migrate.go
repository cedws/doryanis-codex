package db

import (
	"context"
	"database/sql"
	"embed"

	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var migrations embed.FS

func init() {
	goose.SetBaseFS(migrations)

	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}
}

func MigrateUp(ctx context.Context, db *sql.DB) error {
	if err := goose.UpContext(ctx, db, "migrations"); err != nil {
		return err
	}

	return nil
}

func MigrateDown(ctx context.Context, db *sql.DB) error {
	if err := goose.DownToContext(ctx, db, "migrations", 0); err != nil {
		return err
	}

	return nil
}
