package cmd

import (
	"context"
	"encoding/json"
	"os"

	"github.com/pgvector/pgvector-go"
)

type queryCmd struct {
	Query string `arg:""`
}

func (q queryCmd) Run(cli *cli) error {
	ctx := context.Background()

	dbPool, err := connectDB(ctx, cli.DBUsername, cli.DBPassword, cli.DBHost)
	if err != nil {
		return err
	}

	embedding, err := makeEmbeddings(ctx, []string{q.Query})
	if err != nil {
		return err
	}

	vec := pgvector.NewVector(toFloat32Slice(embedding[0]))
	row, err := dbPool.GetMostSimilarActiveSkill(ctx, vec)
	if err != nil {
		return err
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(row)

	return nil
}
