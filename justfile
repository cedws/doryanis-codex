psql:
    PGPASSWORD=postgres psql -h localhost -p 5432 -U postgres -d postgres

generate:
    sqlc generate

format:
    sqlfluff fix --dialect postgres pkg/db
