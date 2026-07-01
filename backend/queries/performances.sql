-- name: CreatePerformance :one
INSERT INTO performances (
    match_id,
    player_id,
    player_team_id,
    title,
    description,
    cover_image_url,
    jersey_number,
    is_starter,
    captain,
    minutes_played,
    provider_rating
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
)
RETURNING id, match_id, player_id, player_team_id, title, description, cover_image_url, jersey_number, is_starter, captain, minutes_played, provider_rating, created_at, updated_at;

-- name: GetPerformanceByID :one
SELECT 
    p.id, p.match_id, p.player_id, p.player_team_id, p.title, p.description, p.cover_image_url, p.jersey_number, p.is_starter, p.captain, p.minutes_played, p.provider_rating, p.created_at, p.updated_at,
    pl.name AS player_name, pl.slug AS player_slug, pl.photo_url AS player_photo_url, pl.known_as AS player_known_as,
    t.name AS team_name, t.short_name AS team_short_name, t.logo_url AS team_logo_url,
    COALESCE(avg_pr.average_rating, 0.0)::numeric(2,1) AS average_rating,
    COALESCE(avg_pr.total_votes, 0)::bigint AS total_votes
FROM performances p
JOIN players pl ON p.player_id = pl.id
JOIN player_teams pt ON p.player_team_id = pt.id
JOIN teams t ON pt.team_id = t.id
LEFT JOIN (
    SELECT performance_id, AVG(rating) AS average_rating, COUNT(id) AS total_votes
    FROM performance_ratings
    GROUP BY performance_id
) avg_pr ON p.id = avg_pr.performance_id
WHERE p.id = $1 LIMIT 1;

-- name: GetPerformanceByPlayer :many
SELECT 
    p.id, p.match_id, p.player_id, p.player_team_id, p.title, p.cover_image_url, p.jersey_number, p.minutes_played, p.provider_rating, p.created_at,
    m.title AS match_title, m.slug AS match_slug, m.utc_datetime AS match_utc_datetime
FROM performances p
JOIN matches m ON p.match_id = m.id
WHERE p.player_id = $1
ORDER BY m.utc_datetime DESC;

-- name: GetPerformanceByMatch :many
SELECT 
    p.id, p.match_id, p.player_id, p.player_team_id, p.title, p.cover_image_url, p.jersey_number, p.is_starter, p.minutes_played, p.provider_rating,
    pl.name AS player_name, pl.photo_url AS player_photo_url
FROM performances p
JOIN players pl ON p.player_id = pl.id
WHERE p.match_id = $1
ORDER BY p.is_starter DESC, p.minutes_played DESC;

-- name: GetPerformanceByPlayerAndMatch :one
SELECT 
    p.id, p.match_id, p.player_id, p.player_team_id, p.title, p.description, p.cover_image_url, p.jersey_number, p.is_starter, p.captain, p.minutes_played, p.provider_rating, p.created_at, p.updated_at
FROM performances p
WHERE p.player_id = $1 AND p.match_id = $2 LIMIT 1;

-- name: SearchPerformances :many
SELECT 
    p.id, p.match_id, p.player_id, p.title, p.cover_image_url, p.provider_rating,
    pl.name AS player_name, m.title AS match_title
FROM performances p
JOIN players pl ON p.player_id = pl.id
JOIN matches m ON p.match_id = m.id
WHERE p.title ILIKE $1 OR pl.name ILIKE $1 OR m.title ILIKE $1
ORDER BY m.utc_datetime DESC
LIMIT $2 OFFSET $3;

-- name: ListPerformances :many
SELECT 
    p.id, p.match_id, p.player_id, p.player_team_id, p.title, p.cover_image_url, p.provider_rating, p.created_at, p.updated_at,
    pl.name AS player_name, m.title AS match_title
FROM performances p
JOIN players pl ON p.player_id = pl.id
JOIN matches m ON p.match_id = m.id
ORDER BY m.utc_datetime DESC
LIMIT $1 OFFSET $2;

-- name: UpdatePerformance :one
UPDATE performances
SET
    match_id = COALESCE($2, match_id),
    player_id = COALESCE($3, player_id),
    player_team_id = COALESCE($4, player_team_id),
    title = COALESCE($5, title),
    description = COALESCE($6, description),
    cover_image_url = COALESCE($7, cover_image_url),
    jersey_number = COALESCE($8, jersey_number),
    is_starter = COALESCE($9, is_starter),
    captain = COALESCE($10, captain),
    minutes_played = COALESCE($11, minutes_played),
    provider_rating = COALESCE($12, provider_rating),
    updated_at = now()
WHERE id = $1
RETURNING id, match_id, player_id, player_team_id, title, description, cover_image_url, jersey_number, is_starter, captain, minutes_played, provider_rating, created_at, updated_at;

-- name: DeletePerformance :exec
DELETE FROM performances
WHERE id = $1;

-- name: GetTopRatedPerformances :many
SELECT 
    p.id, p.match_id, p.player_id, p.title, p.cover_image_url, p.provider_rating,
    pl.name AS player_name, m.title AS match_title,
    COALESCE(avg_pr.average_rating, 0.0)::numeric(2,1) AS average_rating
FROM performances p
JOIN players pl ON p.player_id = pl.id
JOIN matches m ON p.match_id = m.id
JOIN (
    SELECT performance_id, AVG(rating) AS average_rating
    FROM performance_ratings
    GROUP BY performance_id
) avg_pr ON p.id = avg_pr.performance_id
ORDER BY average_rating DESC, m.utc_datetime DESC
LIMIT $1 OFFSET $2;

-- name: GetRecentPerformances :many
SELECT 
    p.id, p.match_id, p.player_id, p.title, p.cover_image_url, p.provider_rating, p.created_at,
    pl.name AS player_name, m.title AS match_title
FROM performances p
JOIN players pl ON p.player_id = pl.id
JOIN matches m ON p.match_id = m.id
ORDER BY p.created_at DESC
LIMIT $1 OFFSET $2;



-- name: GetRelatedPerformances :many
SELECT 
    p.id, p.match_id, p.player_id, p.title, p.cover_image_url, p.provider_rating,
    pl.name AS player_name
FROM performances p
JOIN players pl ON p.player_id = pl.id
WHERE p.match_id = (SELECT sub_p.match_id FROM performances sub_p WHERE sub_p.id = $1) AND p.id <> $1
ORDER BY p.provider_rating DESC
LIMIT $2;
