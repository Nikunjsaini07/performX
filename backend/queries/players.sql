-- name: GetPlayerByID :one
SELECT 
    p.id, p.name, p.slug, p.full_name, p.known_as, p.date_of_birth, p.place_of_birth, p.country_id, p.photo_url, p.height_cm, p.weight_kg, p.shirt_name, p.created_at, p.updated_at,
    c.name AS country_name, c.iso2 AS country_iso2, c.iso3 AS country_iso3
FROM players p
LEFT JOIN countries c ON p.country_id = c.id
WHERE p.id = $1 LIMIT 1;

-- name: GetPlayerBySlug :one
SELECT 
    p.id, p.name, p.slug, p.full_name, p.known_as, p.date_of_birth, p.place_of_birth, p.country_id, p.photo_url, p.height_cm, p.weight_kg, p.shirt_name, p.created_at, p.updated_at,
    c.name AS country_name, c.iso2 AS country_iso2, c.iso3 AS country_iso3
FROM players p
LEFT JOIN countries c ON p.country_id = c.id
WHERE p.slug = $1 LIMIT 1;

-- name: SearchPlayers :many
SELECT 
    p.id, p.name, p.slug, p.known_as, p.photo_url, p.shirt_name,
    c.name AS country_name
FROM players p
LEFT JOIN countries c ON p.country_id = c.id
WHERE p.name ILIKE $1 OR p.known_as ILIKE $1 OR p.full_name ILIKE $1
ORDER BY p.name
LIMIT $2 OFFSET $3;

-- name: ListPlayers :many
SELECT 
    p.id, p.name, p.slug, p.full_name, p.known_as, p.date_of_birth, p.country_id, p.photo_url, p.height_cm, p.weight_kg, p.shirt_name, p.created_at, p.updated_at,
    c.name AS country_name
FROM players p
LEFT JOIN countries c ON p.country_id = c.id
ORDER BY p.name
LIMIT $1 OFFSET $2;

-- name: CreatePlayer :one
INSERT INTO players (
    name,
    slug,
    full_name,
    known_as,
    date_of_birth,
    place_of_birth,
    country_id,
    photo_url,
    height_cm,
    weight_kg,
    shirt_name
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
)
RETURNING id, name, slug, full_name, known_as, date_of_birth, place_of_birth, country_id, photo_url, height_cm, weight_kg, shirt_name, created_at, updated_at;

-- name: UpdatePlayer :one
UPDATE players
SET
    name = COALESCE($2, name),
    slug = COALESCE($3, slug),
    full_name = COALESCE($4, full_name),
    known_as = COALESCE($5, known_as),
    date_of_birth = COALESCE($6, date_of_birth),
    place_of_birth = COALESCE($7, place_of_birth),
    country_id = COALESCE($8, country_id),
    photo_url = COALESCE($9, photo_url),
    height_cm = COALESCE($10, height_cm),
    weight_kg = COALESCE($11, weight_kg),
    shirt_name = COALESCE($12, shirt_name),
    updated_at = now()
WHERE id = $1
RETURNING id, name, slug, full_name, known_as, date_of_birth, place_of_birth, country_id, photo_url, height_cm, weight_kg, shirt_name, created_at, updated_at;

-- name: DeletePlayer :exec
DELETE FROM players
WHERE id = $1;

-- name: GetPlayersByCountry :many
SELECT 
    p.id, p.name, p.slug, p.known_as, p.photo_url, p.shirt_name
FROM players p
WHERE p.country_id = $1
ORDER BY p.name
LIMIT $2 OFFSET $3;

-- name: GetPlayersByTeam :many
SELECT 
    p.id, p.name, p.slug, p.known_as, p.photo_url, p.shirt_name,
    pt.jersey_number, pt.start_date, pt.is_active
FROM player_teams pt
JOIN players p ON pt.player_id = p.id
WHERE pt.team_id = $1 AND pt.is_active = TRUE
ORDER BY p.name;

-- name: GetPlayerPerformances :many
SELECT 
    p.id AS performance_id, p.title AS performance_title, p.cover_image_url AS performance_cover_image, p.minutes_played, p.average_rating, p.created_at,
    m.id AS match_id, m.title AS match_title, m.slug AS match_slug, m.utc_datetime AS match_utc_datetime, m.home_score, m.away_score,
    ht.name AS home_team_name, ht.logo_url AS home_team_logo_url,
    at.name AS away_team_name, at.logo_url AS away_team_logo_url
FROM performances p
JOIN matches m ON p.match_id = m.id
JOIN teams ht ON m.home_team_id = ht.id
JOIN teams at ON m.away_team_id = at.id
WHERE p.player_id = $1
ORDER BY m.utc_datetime DESC
LIMIT $2 OFFSET $3;

-- name: GetPlayerTopRatedPerformances :many
SELECT 
    p.id AS performance_id, p.title AS performance_title, p.cover_image_url AS performance_cover_image, p.minutes_played, p.average_rating,
    m.id AS match_id, m.title AS match_title, m.slug AS match_slug, m.utc_datetime AS match_utc_datetime
FROM performances p
JOIN matches m ON p.match_id = m.id
WHERE p.player_id = $1
ORDER BY p.average_rating DESC, m.utc_datetime DESC
LIMIT $2 OFFSET $3;

-- name: GetPlayerAverageRating :one
SELECT 
    COALESCE(AVG(pr.rating), 0.0)::numeric(2,1) AS average_rating,
    COUNT(pr.id)::bigint AS total_votes
FROM performances p
JOIN performance_ratings pr ON p.id = pr.performance_id
WHERE p.player_id = $1;

-- name: GetPlayerReviews :many
SELECT 
    r.id AS review_id, r.title AS review_title, r.content AS review_content, r.created_at AS review_created_at,
    u.id AS user_id, u.username, u.display_name, u.avatar_url,
    p.id AS performance_id, p.title AS performance_title,
    m.id AS match_id, m.title AS match_title
FROM performance_reviews r
JOIN performances p ON r.performance_id = p.id
JOIN matches m ON p.match_id = m.id
JOIN users u ON r.user_id = u.id
WHERE p.player_id = $1
ORDER BY r.created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetPlayerReviewCount :one
SELECT COUNT(r.id)::bigint AS review_count
FROM performance_reviews r
JOIN performances p ON r.performance_id = p.id
WHERE p.player_id = $1;

-- name: GetPlayerMatches :many
SELECT 
    m.id, m.title, m.slug, m.utc_datetime, m.home_score, m.away_score,
    ht.name AS home_team_name, ht.logo_url AS home_team_logo_url,
    at.name AS away_team_name, at.logo_url AS away_team_logo_url
FROM performances p
JOIN matches m ON p.match_id = m.id
JOIN teams ht ON m.home_team_id = ht.id
JOIN teams at ON m.away_team_id = at.id
WHERE p.player_id = $1
ORDER BY m.utc_datetime DESC
LIMIT $2 OFFSET $3;

-- name: GetRelatedPlayers :many
SELECT DISTINCT 
    p.id, p.name, p.slug, p.known_as, p.photo_url, p.shirt_name
FROM players p
LEFT JOIN player_teams pt ON p.id = pt.player_id
WHERE p.id <> $1 AND (
    p.country_id = (SELECT country_id FROM players WHERE id = $1)
    OR pt.team_id IN (SELECT team_id FROM player_teams WHERE player_id = $1 AND is_active = TRUE)
)
LIMIT $2;

-- name: GetPlayerAggregatedStats :many
SELECT 
    st.id AS stat_type_id,
    st.name AS stat_name,
    st.short_name AS stat_short_name,
    st.unit AS stat_unit,
    st.category AS stat_category,
    SUM(ps.value)::double precision AS total_value,
    AVG(ps.value)::double precision AS average_value
FROM performance_stats ps
JOIN stat_types st ON ps.stat_type_id = st.id
JOIN performances p ON ps.performance_id = p.id
WHERE p.player_id = $1
GROUP BY st.id, st.name, st.short_name, st.unit, st.category, st.display_order
ORDER BY st.category, st.display_order;

-- name: GetPlayerRatings :many
SELECT 
    pr.id AS rating_id, pr.rating, pr.created_at AS rated_at,
    u.id AS user_id, u.username, u.display_name, u.avatar_url,
    p.id AS performance_id, p.title AS performance_title,
    m.id AS match_id, m.title AS match_title
FROM performance_ratings pr
JOIN performances p ON pr.performance_id = p.id
JOIN matches m ON p.match_id = m.id
JOIN users u ON r.user_id = u.id
WHERE p.player_id = $1
ORDER BY pr.created_at DESC
LIMIT $2 OFFSET $3;
