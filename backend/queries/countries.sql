-- name: CreateCountry :one
INSERT INTO countries (
    name,
    iso2,
    iso3
) VALUES (
    $1, $2, $3
)
RETURNING id, name, iso2, iso3, created_at, updated_at;

-- name: GetCountryByID :one
SELECT id, name, iso2, iso3, created_at, updated_at
FROM countries
WHERE id = $1 LIMIT 1;

-- name: GetCountryByCode :one
SELECT id, name, iso2, iso3, created_at, updated_at
FROM countries
WHERE iso2 = UPPER($1) OR iso3 = UPPER($1) LIMIT 1;

-- name: ListCountries :many
SELECT id, name, iso2, iso3, created_at, updated_at
FROM countries
ORDER BY name;

-- name: SearchCountries :many
SELECT id, name, iso2, iso3, created_at, updated_at
FROM countries
WHERE name ILIKE $1 OR iso2 ILIKE $1 OR iso3 ILIKE $1
ORDER BY name;

-- name: UpdateCountry :one
UPDATE countries
SET 
    name = COALESCE($2, name),
    iso2 = COALESCE($3, iso2),
    iso3 = COALESCE($4, iso3),
    updated_at = now()
WHERE id = $1
RETURNING id, name, iso2, iso3, created_at, updated_at;

-- name: DeleteCountry :exec
DELETE FROM countries
WHERE id = $1;
