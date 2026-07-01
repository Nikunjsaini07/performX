-- name: CreateSport :one
INSERT INTO sports (
    name,
    slug
) VALUES (
    $1, $2
)
RETURNING id, name, slug, created_at;

-- name: GetSportByID :one
SELECT id, name, slug, created_at
FROM sports
WHERE id = $1 LIMIT 1;

-- name: GetSportBySlug :one
SELECT id, name, slug, created_at
FROM sports
WHERE slug = $1 LIMIT 1;

-- name: ListSports :many
SELECT id, name, slug, created_at
FROM sports
ORDER BY name;

-- name: UpdateSport :one
UPDATE sports
SET 
    name = COALESCE($2, name),
    slug = COALESCE($3, slug)
WHERE id = $1
RETURNING id, name, slug, created_at;

-- name: DeleteSport :exec
DELETE FROM sports
WHERE id = $1;
