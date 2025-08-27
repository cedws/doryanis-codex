package codex

import (
	"context"
	"encoding/json"
	"log/slog"
	"maps"
	"os"
	"slices"

	"github.com/cedws/doryanis-codex/pkg/db"
	"github.com/cedws/doryanis-codex/pkg/types"
	"github.com/pgvector/pgvector-go"
)

func LoadData(ctx context.Context, opts Options, dataPath string) error {
	_, dbQuery, err := connectDB(ctx, opts.DBUsername, opts.DBPassword, opts.DBHost)
	if err != nil {
		return err
	}

	file, err := loadDataFile(dataPath)
	if err != nil {
		return err
	}

	skills := collectSkills(file)

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
	skillChunks := slices.Chunk(skills, 10)

	embedClient := NewEmbeddingsClient[int]()

	for skillChunk := range skillChunks {
		embedInputs := make(EmbeddingsBatchRequest[int])

		for i, skill := range skillChunk {
			text, err := skill.MarshalText()
			if err != nil {
				return err
			}
			embedInputs[i] = string(text)
		}

		embedResults, err := embedClient.BatchEmbed(ctx, embedInputs)
		if err != nil {
			return err
		}

		if err := createActiveSkills(ctx, dbPool, skillChunk, embedResults); err != nil {
			return err
		}
	}

	return nil
}

func createActiveSkills(ctx context.Context, dbPool *db.Queries, skills []types.ActiveSkill, embeddings EmbeddingsBatchResult[int]) error {
	for i, skill := range skills {
		vec := pgvector.NewVector(embeddings[i])

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

		slog.Info("processed active skill", "name", skill.DisplayName)
	}
	return nil
}
