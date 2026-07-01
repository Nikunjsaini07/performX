-- name: CreateTournament :one
INSERT INTO tournaments (
    sport_id,
    name,
    short_name,
    slug,
    type,
    logo_url
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING id, sport_id, name, short_name, slug, type, logo_url, created_at, updated_at;

-- name: GetTournamentByID :one
SELECT id, sport_id, name, short_name, slug, type, logo_url, created_at, updated_at
FROM tournaments
WHERE id = $1 LIMIT 1;

-- name: GetTournamentBySlug :one
SELECT id, sport_id, name, short_name, slug, type, logo_url, created_at, updated_at
FROM tournaments
WHERE slug = $1 LIMIT 1;

-- name: ListTournaments :many
SELECT id, sport_id, name, short_name, slug, type, logo_url, created_at, updated_at
FROM tournaments
ORDER BY name;

-- name: GetTournamentsBySport :many
SELECT id, sport_id, name, short_name, slug, type, logo_url, created_at, updated_at
FROM tournaments
WHERE sport_id = $1
ORDER BY name;

-- name: UpdateTournament :one
UPDATE tournaments
SET 
    sport_id = COALESCE($2, sport_id),
    name = COALESCE($3, name),
    short_name = COALESCE($4, short_name),
    slug = COALESCE($5, slug),
    type = COALESCE($6, type),
    logo_url = COALESCE($7, logo_url),
    updated_at = now()
WHERE id = $1
RETURNING id, sport_id, name, short_name, slug, type, logo_url, created_at, updated_at;

-- name: DeleteTournament :exec
DELETE FROM tournaments
WHERE id = $1;
