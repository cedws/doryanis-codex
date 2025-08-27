package cmd

import (
	"github.com/alecthomas/kong"
)

type cli struct {
	DBUsername string
	DBPassword string
	DBHost     string

	Query    queryCmd    `cmd:""`
	LoadData loadDataCmd `cmd:""`
	Migrate  migrateCmd  `cmd:""`
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
