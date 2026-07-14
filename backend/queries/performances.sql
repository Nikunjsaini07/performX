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
    average_rating,
    slug,
    tagline
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
)
RETURNING id, match_id, player_id, player_team_id, title, description, cover_image_url, jersey_number, is_starter, captain, minutes_played, average_rating, total_votes, slug, tagline, created_at, updated_at;

-- name: GetPerformanceByID :one
SELECT 
    p.id, p.match_id, p.player_id, p.player_team_id, p.title, p.description, p.cover_image_url, p.jersey_number, p.is_starter, p.captain, p.minutes_played, p.average_rating, p.slug, p.tagline, p.created_at, p.updated_at,
    pl.name AS player_name, pl.slug AS player_slug, pl.photo_url AS player_photo_url, pl.known_as AS player_known_as,
    t.name AS team_name, t.short_name AS team_short_name, t.logo_url AS team_logo_url,
    p.total_votes
FROM performances p
JOIN players pl ON p.player_id = pl.id
JOIN player_teams pt ON p.player_team_id = pt.id
JOIN teams t ON pt.team_id = t.id
WHERE p.id = $1 LIMIT 1;

-- name: GetPerformanceBySlug :one
SELECT 
    p.id, p.match_id, p.player_id, p.player_team_id, p.title, p.description, p.cover_image_url, p.jersey_number, p.is_starter, p.captain, p.minutes_played, p.average_rating, p.slug, p.tagline, p.created_at, p.updated_at,
    pl.name AS player_name, pl.slug AS player_slug, pl.photo_url AS player_photo_url, pl.known_as AS player_known_as,
    t.name AS team_name, t.short_name AS team_short_name, t.logo_url AS team_logo_url,
    p.total_votes
FROM performances p
JOIN players pl ON p.player_id = pl.id
JOIN player_teams pt ON p.player_team_id = pt.id
JOIN teams t ON pt.team_id = t.id
WHERE p.slug = $1 LIMIT 1;

-- name: GetPerformanceByPlayer :many
SELECT 
    p.id, p.match_id, p.player_id, p.player_team_id, p.title, p.cover_image_url, p.jersey_number, p.minutes_played, p.average_rating, p.created_at,
    m.title AS match_title, m.slug AS match_slug, m.utc_datetime AS match_utc_datetime
FROM performances p
JOIN matches m ON p.match_id = m.id
WHERE p.player_id = $1
ORDER BY m.utc_datetime DESC;

-- name: GetPerformanceByMatch :many
SELECT 
    p.id, p.match_id, p.player_id, p.player_team_id, p.title, p.cover_image_url, p.jersey_number, p.is_starter, p.minutes_played, p.average_rating, p.slug, p.tagline,
    pl.name AS player_name, pl.photo_url AS player_photo_url,
    'N/A' AS position, t.name AS team_name, t.flag_emoji,
    COALESCE((SELECT sv.value FROM performance_stats sv JOIN stat_types stt ON stt.id = sv.stat_type_id WHERE sv.performance_id = p.id AND stt.name = 'goals' LIMIT 1), 0)::int AS goals,
    COALESCE((SELECT sv.value FROM performance_stats sv JOIN stat_types stt ON stt.id = sv.stat_type_id WHERE sv.performance_id = p.id AND stt.name = 'assists' LIMIT 1), 0)::int AS assists,
    0::int AS dribbles,
    0::numeric(5,2) AS passes_accuracy
FROM performances p
JOIN players pl ON p.player_id = pl.id
LEFT JOIN player_teams pt ON p.player_team_id = pt.id
LEFT JOIN teams t ON pt.team_id = t.id
WHERE p.match_id = $1
ORDER BY p.is_starter DESC, p.minutes_played DESC;

-- name: GetPerformanceByPlayerAndMatch :one
SELECT 
    p.id, p.match_id, p.player_id, p.player_team_id, p.title, p.description, p.cover_image_url, p.jersey_number, p.is_starter, p.captain, p.minutes_played, p.average_rating, p.slug, p.tagline, p.created_at, p.updated_at
FROM performances p
WHERE p.player_id = $1 AND p.match_id = $2 LIMIT 1;

-- name: SearchPerformances :many
SELECT 
    p.id, p.match_id, p.player_id, p.title, p.cover_image_url, p.average_rating,
    pl.name AS player_name, m.title AS match_title
FROM performances p
JOIN players pl ON p.player_id = pl.id
JOIN matches m ON p.match_id = m.id
WHERE p.title ILIKE $1 OR pl.name ILIKE $1 OR m.title ILIKE $1
ORDER BY m.utc_datetime DESC
LIMIT $2 OFFSET $3;

-- name: ListPerformances :many
SELECT 
    p.id, p.match_id, p.player_id, p.player_team_id, p.title, p.cover_image_url, p.average_rating, p.created_at, p.updated_at, p.minutes_played, p.slug, p.tagline,
    pl.name AS player_name, m.title AS match_title, m.round AS match_round,
    'N/A' AS position, t.name AS team_name, t.flag_emoji,
    COALESCE((SELECT sv.value FROM performance_stats sv JOIN stat_types stt ON stt.id = sv.stat_type_id WHERE sv.performance_id = p.id AND stt.name = 'goals' LIMIT 1), 0)::int AS goals,
    COALESCE((SELECT sv.value FROM performance_stats sv JOIN stat_types stt ON stt.id = sv.stat_type_id WHERE sv.performance_id = p.id AND stt.name = 'assists' LIMIT 1), 0)::int AS assists,
    0::int AS dribbles,
    0::numeric(5,2) AS passes_accuracy,
    p.total_votes
FROM performances p
JOIN players pl ON p.player_id = pl.id
JOIN matches m ON p.match_id = m.id
LEFT JOIN player_teams pt ON p.player_team_id = pt.id
LEFT JOIN teams t ON pt.team_id = t.id
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
    average_rating = COALESCE($12, average_rating),
    slug = COALESCE($13, slug),
    tagline = COALESCE($14, tagline),
    updated_at = now()
WHERE id = $1
RETURNING id, match_id, player_id, player_team_id, title, description, cover_image_url, jersey_number, is_starter, captain, minutes_played, average_rating, slug, tagline, created_at, updated_at;

-- name: DeletePerformance :exec
DELETE FROM performances
WHERE id = $1;

-- name: GetTopRatedPerformances :many
SELECT 
    p.id, p.match_id, p.player_id, p.title, p.cover_image_url, p.average_rating, p.total_votes,
    pl.name AS player_name, m.title AS match_title
FROM performances p
JOIN players pl ON p.player_id = pl.id
JOIN matches m ON p.match_id = m.id
ORDER BY p.average_rating DESC, m.utc_datetime DESC
LIMIT $1 OFFSET $2;

-- name: GetRecentPerformances :many
SELECT 
    p.id, p.match_id, p.player_id, p.title, p.cover_image_url, p.average_rating, p.created_at,
    pl.name AS player_name, m.title AS match_title
FROM performances p
JOIN players pl ON p.player_id = pl.id
JOIN matches m ON p.match_id = m.id
ORDER BY p.created_at DESC
LIMIT $1 OFFSET $2;



-- name: GetRelatedPerformances :many
SELECT 
    p.id, p.match_id, p.player_id, p.title, p.cover_image_url, p.average_rating,
    pl.name AS player_name
FROM performances p
JOIN players pl ON p.player_id = pl.id
WHERE p.match_id = (SELECT sub_p.match_id FROM performances sub_p WHERE sub_p.id = $1) AND p.id <> $1
ORDER BY p.average_rating DESC
LIMIT $2;
