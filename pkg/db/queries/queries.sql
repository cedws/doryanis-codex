-- name: CreateActiveSkill :one
INSERT INTO active_skills (
    display_name,
    description,
    types,
    embedding_id
) VALUES (
    sqlc.arg(display_name)::text,
    sqlc.arg(description)::text,
    sqlc.arg(types),
    sqlc.arg(embedding_id)::bigint
)
RETURNING
    id;

-- name: CreateEmbedding :one
INSERT INTO embeddings (
    embedding
) VALUES (
    sqlc.arg(embedding)::vector(3072)
)
RETURNING
    id;

-- name: GetMostSimilarActiveSkills :many
SELECT
    a.id,
    a.display_name,
    a.description,
    a.types
FROM active_skills AS a
INNER JOIN embeddings AS e ON a.embedding_id = e.id
WHERE a.embedding_id IS NOT NULL
ORDER BY e.embedding <-> sqlc.arg(query_embedding)::vector(3072)
LIMIT sqlc.arg(n)::int;
