package cmd

import (
	"context"
	"fmt"

	"github.com/cedws/doryanis-codex/pkg/db"
	"github.com/jackc/pgx/v5/stdlib"
)

type migrateCmd struct {
	Up   migrateUpCmd   `cmd:""`
	Down migrateDownCmd `cmd:""`
}

type migrateUpCmd struct{}

func (m migrateUpCmd) Run(cli *cli) error {
	ctx := context.Background()

	dbPool, _, err := connectDB(ctx, cli.DBUsername, cli.DBPassword, cli.DBHost)
	if err != nil {
		return err
	}

	if err := db.MigrateUp(ctx, stdlib.OpenDBFromPool(dbPool)); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	return nil
}

type migrateDownCmd struct{}

func (m migrateDownCmd) Run(cli *cli) error {
	ctx := context.Background()

	dbPool, _, err := connectDB(ctx, cli.DBUsername, cli.DBPassword, cli.DBHost)
	if err != nil {
		return err
	}

	if err := db.MigrateDown(ctx, stdlib.OpenDBFromPool(dbPool)); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	return nil
}
