package cmd

import (
	"context"
	"encoding"
	"encoding/json"
	"fmt"
	"log"
	"maps"
	"os"
	"slices"

	"github.com/cedws/doryanis-codex/pkg/db"
	"github.com/cedws/doryanis-codex/pkg/types"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/openai/openai-go"
	"github.com/pgvector/pgvector-go"
)

type loadDataCmd struct {
	DataPath string
}

func (l loadDataCmd) Run(cli *cli) error {
	ctx := context.Background()

	dbPool, dbQuery, err := connectDB(ctx, cli.DBUsername, cli.DBPassword, cli.DBHost)
	if err != nil {
		return err
	}

	if err := db.Migrate(ctx, stdlib.OpenDBFromPool(dbPool)); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	data, err := loadDataFile(l.DataPath)
	if err != nil {
		return err
	}

	skills := collectSkills(data)

	return insertSkills(ctx, dbQuery, skills)
}

func loadDataFile(path string) (*types.Data, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var data types.Data
	if err := json.NewDecoder(f).Decode(&data); err != nil {
		return nil, err
	}

	return &data, nil
}

func collectSkills(data *types.Data) []types.ActiveSkill {
	return slices.Collect(maps.Values(data.ActiveSkills))
}

func insertSkills(ctx context.Context, dbPool *db.Queries, skills []types.ActiveSkill) error {
	chunks := slices.Chunk(skills, 10)
	openaiClient := openai.NewClient()

	for chunk := range chunks {
		ents := make([]encoding.TextMarshaler, 0, len(chunk))
		for _, s := range chunk {
			ents = append(ents, s)
		}

		embeddings, err := makeEntEmbeddings(ctx, openaiClient, ents)
		if err != nil {
			return err
		}

		if err := saveChunk(ctx, dbPool, chunk, embeddings); err != nil {
			return err
		}
	}

	return nil
}

func saveChunk(ctx context.Context, dbPool *db.Queries, skills []types.ActiveSkill, embeddings [][]float64) error {
	for i, skill := range skills {
		vec := pgvector.NewVector(toFloat32Slice(embeddings[i]))

		id, err := dbPool.CreateEmbedding(ctx, vec)
		if err != nil {
			return err
		}

		_, err = dbPool.CreateActiveSkill(ctx, db.CreateActiveSkillParams{
			DisplayName: skill.DisplayName,
			Description: skill.Description,
			Types:       skill.Types,
			EmbeddingID: id,
		})
		if err != nil {
			return err
		}

		log.Printf("Processed %s\n", skill.DisplayName)
	}
	return nil
}

func toFloat32Slice(f []float64) []float32 {
	out := make([]float32, len(f))
	for i, v := range f {
		out[i] = float32(v)
	}
	return out
}

func makeEntEmbeddings(ctx context.Context, client openai.Client, entities []encoding.TextMarshaler) ([][]float64, error) {
	var input []string
	for _, ent := range entities {
		text, err := ent.MarshalText()
		if err != nil {
			return nil, err
		}
		input = append(input, string(text))
	}
	return makeEmbeddings(ctx, client, input)
}

func makeEmbeddings(ctx context.Context, client openai.Client, input []string) ([][]float64, error) {
	resp, err := client.Embeddings.New(ctx, openai.EmbeddingNewParams{
		Input: openai.EmbeddingNewParamsInputUnion{
			OfArrayOfStrings: input,
		},
		Model: openai.EmbeddingModelTextEmbedding3Large,
	})
	if err != nil {
		return nil, err
	}

	output := make([][]float64, len(input))
	for _, data := range resp.Data {
		output[data.Index] = data.Embedding
	}

	return output, nil
}
