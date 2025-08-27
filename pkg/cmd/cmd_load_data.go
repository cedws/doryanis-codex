package cmd

import (
	"context"

	"github.com/cedws/doryanis-codex/pkg/codex"
)

type loadDataCmd struct {
	DataPath string
}

func (l loadDataCmd) Run(cli *cli) error {
	ctx := context.Background()

	opts := codex.Options{
		DBUsername: cli.DBUsername,
		DBPassword: cli.DBPassword,
		DBHost:     cli.DBHost,
	}

	return codex.LoadData(ctx, opts, l.DataPath)
}
