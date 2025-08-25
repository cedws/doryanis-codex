package cmd

import (
	"context"
	"encoding/json"
	"os"

	"github.com/openai/openai-go"
	"github.com/pgvector/pgvector-go"
)

type queryCmd struct {
	Query string `arg:""`
}

func (q queryCmd) Run(cli *cli) error {
	ctx := context.Background()

	_, dbQuery, err := connectDB(ctx, cli.DBUsername, cli.DBPassword, cli.DBHost)
	if err != nil {
		return err
	}

	openaiClient := openai.NewClient()

	embedding, err := makeEmbeddings(ctx, openaiClient, []string{q.Query})
	if err != nil {
		return err
	}

	vec := pgvector.NewVector(toFloat32Slice(embedding[0]))

	row, err := dbQuery.GetMostSimilarActiveSkill(ctx, vec)
	if err != nil {
		return err
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(row)

	return nil
}
