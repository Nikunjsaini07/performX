-- name: CreateSeason :one
INSERT INTO seasons (
    tournament_id,
    name,
    start_year,
    end_year,
    start_date,
    end_date,
    is_current
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING id, tournament_id, name, start_year, end_year, start_date, end_date, is_current, created_at, updated_at;

-- name: GetSeasonByID :one
SELECT id, tournament_id, name, start_year, end_year, start_date, end_date, is_current, created_at, updated_at
FROM seasons
WHERE id = $1 LIMIT 1;

-- name: GetSeasonBySlug :one
SELECT id, tournament_id, name, start_year, end_year, start_date, end_date, is_current, created_at, updated_at
FROM seasons
WHERE lower(replace(name, ' ', '-')) = $1 LIMIT 1;

-- name: ListSeasons :many
SELECT id, tournament_id, name, start_year, end_year, start_date, end_date, is_current, created_at, updated_at
FROM seasons
ORDER BY start_year DESC;

-- name: ListSeasonsByTournament :many
SELECT id, tournament_id, name, start_year, end_year, start_date, end_date, is_current, created_at, updated_at
FROM seasons
WHERE tournament_id = $1
ORDER BY start_year DESC;

-- name: UpdateSeason :one
UPDATE seasons
SET 
    tournament_id = COALESCE($2, tournament_id),
    name = COALESCE($3, name),
    start_year = COALESCE($4, start_year),
    end_year = COALESCE($5, end_year),
    start_date = COALESCE($6, start_date),
    end_date = COALESCE($7, end_date),
    is_current = COALESCE($8, is_current),
    updated_at = now()
WHERE id = $1
RETURNING id, tournament_id, name, start_year, end_year, start_date, end_date, is_current, created_at, updated_at;

-- name: DeleteSeason :exec
DELETE FROM seasons
WHERE id = $1;
