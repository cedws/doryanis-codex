-- +goose Up
BEGIN;

CREATE EXTENSION IF NOT EXISTS vector;

CREATE TABLE embeddings (
  id BIGSERIAL PRIMARY KEY,
  embedding vector(3072) NOT NULL
);

CREATE TABLE active_skills (
  id BIGSERIAL PRIMARY KEY,
  display_name TEXT,
  description TEXT,
  types JSONB,
  embedding_id BIGINT REFERENCES embeddings(id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

COMMIT;
-- +goose Down
BEGIN;
DROP TABLE IF EXISTS active_skills;
DROP TABLE IF EXISTS embeddings;
COMMIT;
