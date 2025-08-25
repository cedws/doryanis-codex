-- name: CreateActiveSkill :one
INSERT INTO active_skills (
  display_name,
  description,
  types,
  embedding
) VALUES (
  sqlc.arg(display_name)::text,
  sqlc.arg(description)::text,
  sqlc.arg(types),
  sqlc.arg(embedding)::vector(3072)
)
RETURNING
  id;

-- name: GetMostSimilarActiveSkill :one
SELECT
  id,
  display_name,
  description,
  types
FROM active_skills
ORDER BY embedding <-> sqlc.arg(query_embedding)::vector(3072)
LIMIT 1;
