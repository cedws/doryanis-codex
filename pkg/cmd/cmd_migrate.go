package cmd

import (
	"context"

	"github.com/cedws/doryanis-codex/pkg/codex"
)

type migrateCmd struct {
	Up   migrateUpCmd   `cmd:""`
	Down migrateDownCmd `cmd:""`
}

type migrateUpCmd struct{}

func (m migrateUpCmd) Run(cli *cli) error {
	ctx := context.Background()

	return codex.MigrateUp(ctx, codex.Options{
		DBUsername: cli.DBUsername,
		DBPassword: cli.DBPassword,
		DBHost:     cli.DBHost,
	})
}

type migrateDownCmd struct{}

func (m migrateDownCmd) Run(cli *cli) error {
	ctx := context.Background()

	return codex.MigrateDown(ctx, codex.Options{
		DBUsername: cli.DBUsername,
		DBPassword: cli.DBPassword,
		DBHost:     cli.DBHost,
	})
}
