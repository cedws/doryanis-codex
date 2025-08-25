package cmd

import (
	"context"
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/cedws/doryanis-codex/pkg/db"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

func connectDB(ctx context.Context, username, password, host string) (*db.Queries, error) {
	dbPool, err := pgxpool.New(context.Background(), fmt.Sprintf("postgres://%s:%s@%s", username, password, host))
	if err != nil {
		return nil, err
	}

	if err := dbPool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	sqlDB := stdlib.OpenDBFromPool(dbPool)

	if err := db.Migrate(ctx, sqlDB); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return db.New(dbPool), nil
}

type cli struct {
	DBUsername string
	DBPassword string
	DBHost     string

	Query    queryCmd    `cmd:""`
	LoadData loadDataCmd `cmd:""`
}

func Execute() {
	var cli cli

	ctx := kong.Parse(&cli,
		kong.Name("doryanis-codex"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
		}))

	ctx.FatalIfErrorf(ctx.Run())
}
