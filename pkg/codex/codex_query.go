package codex

import (
	"context"
	"encoding/json"
	"os"

	"github.com/cedws/doryanis-codex/pkg/db"
	"github.com/pgvector/pgvector-go"
)

func Query(ctx context.Context, opts Options, query string) error {
	_, dbQuery, err := connectDB(ctx, opts.DBUsername, opts.DBPassword, opts.DBHost)
	if err != nil {
		return err
	}

	embedClient := NewEmbeddingsClient[string]()
	embedResult, err := embedClient.BatchEmbed(ctx, EmbeddingsBatchRequest[string]{
		"query": query,
	})
	if err != nil {
		return err
	}

	vec := pgvector.NewVector(embedResult["query"])

	row, err := dbQuery.GetMostSimilarActiveSkills(ctx, db.GetMostSimilarActiveSkillsParams{
		N:              5,
		QueryEmbedding: vec,
	})
	if err != nil {
		return err
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(row)

	return nil
}
