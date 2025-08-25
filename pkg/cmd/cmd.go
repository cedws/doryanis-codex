package cmd

import (
	"context"
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/cedws/doryanis-codex/pkg/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

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
