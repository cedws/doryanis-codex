-- +goose Up
BEGIN;

CREATE EXTENSION IF NOT EXISTS vector;

CREATE TABLE active_skills (
  id BIGSERIAL PRIMARY KEY,
  display_name TEXT,
  description TEXT,
  types JSONB,
  embedding vector(3072),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

COMMIT;
-- +goose Down
BEGIN;
DROP TABLE IF EXISTS active_skills;
COMMIT;
