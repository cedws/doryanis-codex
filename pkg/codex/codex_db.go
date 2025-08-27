package codex

import (
	"context"
	"fmt"

	"github.com/cedws/doryanis-codex/pkg/db"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

type Options struct {
	DBUsername string
	DBPassword string
	DBHost     string
}

func connectDB(ctx context.Context, username, password, host string) (*pgxpool.Pool, *db.Queries, error) {
	dbPool, err := pgxpool.New(context.Background(), fmt.Sprintf("postgres://%s:%s@%s", username, password, host))
	if err != nil {
		return nil, nil, err
	}

	if err := dbPool.Ping(ctx); err != nil {
		return nil, nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return dbPool, db.New(dbPool), nil
}

func MigrateUp(ctx context.Context, opts Options) error {
	dbPool, _, err := connectDB(ctx, opts.DBUsername, opts.DBPassword, opts.DBHost)
	if err != nil {
		return err
	}

	if err := db.MigrateUp(ctx, stdlib.OpenDBFromPool(dbPool)); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	return nil
}

func MigrateDown(ctx context.Context, opts Options) error {
	dbPool, _, err := connectDB(ctx, opts.DBUsername, opts.DBPassword, opts.DBHost)
	if err != nil {
		return err
	}

	if err := db.MigrateDown(ctx, stdlib.OpenDBFromPool(dbPool)); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	return nil
}
