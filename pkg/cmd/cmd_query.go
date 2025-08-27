package cmd

import (
	"context"

	"github.com/cedws/doryanis-codex/pkg/codex"
)

type queryCmd struct {
	Query string `arg:""`
}

func (q queryCmd) Run(cli *cli) error {
	ctx := context.Background()

	opts := codex.Options{
		DBUsername: cli.DBUsername,
		DBPassword: cli.DBPassword,
		DBHost:     cli.DBHost,
	}

	return codex.Query(ctx, opts, q.Query)
}
