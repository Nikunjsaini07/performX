-- name: CreateMatch :one
INSERT INTO matches (
    home_team_id,
    away_team_id,
    title,
    slug,
    description,
    round,
    utc_datetime,
    venue,
    cover_image_url,
    home_score,
    away_score,
    home_penalty_score,
    away_penalty_score,
    tagline
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
)
RETURNING id, home_team_id, away_team_id, title, slug, description, round, utc_datetime, venue, cover_image_url, home_score, away_score, home_penalty_score, away_penalty_score, tagline, created_at, updated_at;

-- name: GetMatchByID :one
SELECT 
    m.id, m.home_team_id, m.away_team_id, m.title, m.slug, m.description, m.round, m.utc_datetime, m.venue, m.cover_image_url, m.home_score, m.away_score, m.home_penalty_score, m.away_penalty_score, m.tagline, m.created_at, m.updated_at,
    ht.name AS home_team_name, ht.short_name AS home_team_short_name, ht.logo_url AS home_team_logo_url,
    at.name AS away_team_name, at.short_name AS away_team_short_name, at.logo_url AS away_team_logo_url,
    m.average_rating,
    m.total_votes
FROM matches m
JOIN teams ht ON m.home_team_id = ht.id
JOIN teams at ON m.away_team_id = at.id
WHERE m.id = $1 LIMIT 1;

-- name: GetMatchBySlug :one
SELECT 
    m.id, m.home_team_id, m.away_team_id, m.title, m.slug, m.description, m.round, m.utc_datetime, m.venue, m.cover_image_url, m.home_score, m.away_score, m.home_penalty_score, m.away_penalty_score, m.tagline, m.created_at, m.updated_at,
    ht.name AS home_team_name, ht.short_name AS home_team_short_name, ht.logo_url AS home_team_logo_url,
    at.name AS away_team_name, at.short_name AS away_team_short_name, at.logo_url AS away_team_logo_url,
    m.average_rating,
    m.total_votes
FROM matches m
JOIN teams ht ON m.home_team_id = ht.id
JOIN teams at ON m.away_team_id = at.id
WHERE m.slug = $1 LIMIT 1;

-- name: SearchMatches :many
SELECT 
    m.id, m.home_team_id, m.away_team_id, m.title, m.slug, m.tagline, m.utc_datetime, m.venue, m.home_score, m.away_score,
    ht.name AS home_team_name, ht.logo_url AS home_team_logo_url,
    at.name AS away_team_name, at.logo_url AS away_team_logo_url
FROM matches m
JOIN teams ht ON m.home_team_id = ht.id
JOIN teams at ON m.away_team_id = at.id
WHERE m.title ILIKE $1 OR ht.name ILIKE $1 OR at.name ILIKE $1
ORDER BY m.utc_datetime DESC
LIMIT $2 OFFSET $3;

-- name: ListMatches :many
SELECT 
    m.id, m.home_team_id, m.away_team_id, m.title, m.slug, m.description, m.round, m.utc_datetime, m.venue, m.cover_image_url, m.home_score, m.away_score, m.home_penalty_score, m.away_penalty_score, m.tagline, m.created_at, m.updated_at,
    ht.name AS home_team_name, ht.short_name AS home_team_short_name, ht.logo_url AS home_team_logo_url,
    at.name AS away_team_name, at.short_name AS away_team_short_name, at.logo_url AS away_team_logo_url,
    m.average_rating,
    m.total_votes
FROM matches m
JOIN teams ht ON m.home_team_id = ht.id
JOIN teams at ON m.away_team_id = at.id
ORDER BY m.utc_datetime DESC
LIMIT $1 OFFSET $2;

-- name: GetMatchesByTeam :many
SELECT 
    m.id, m.home_team_id, m.away_team_id, m.title, m.slug, m.tagline, m.utc_datetime, m.home_score, m.away_score,
    ht.name AS home_team_name, ht.logo_url AS home_team_logo_url,
    at.name AS away_team_name, at.logo_url AS away_team_logo_url
FROM matches m
JOIN teams ht ON m.home_team_id = ht.id
JOIN teams at ON m.away_team_id = at.id
WHERE m.home_team_id = $1 OR m.away_team_id = $1
ORDER BY m.utc_datetime DESC
LIMIT $2 OFFSET $3;

-- name: GetMatchesByDate :many
SELECT 
    m.id, m.home_team_id, m.away_team_id, m.title, m.slug, m.tagline, m.utc_datetime, m.home_score, m.away_score,
    ht.name AS home_team_name, ht.logo_url AS home_team_logo_url,
    at.name AS away_team_name, at.logo_url AS away_team_logo_url
FROM matches m
JOIN teams ht ON m.home_team_id = ht.id
JOIN teams at ON m.away_team_id = at.id
WHERE m.utc_datetime::date = $1::date
ORDER BY m.utc_datetime DESC;

-- name: GetUpcomingMatches :many
SELECT 
    m.id, m.home_team_id, m.away_team_id, m.title, m.slug, m.tagline, m.utc_datetime, m.venue,
    ht.name AS home_team_name, ht.logo_url AS home_team_logo_url,
    at.name AS away_team_name, at.logo_url AS away_team_logo_url
FROM matches m
JOIN teams ht ON m.home_team_id = ht.id
JOIN teams at ON m.away_team_id = at.id
WHERE m.utc_datetime > now()
ORDER BY m.utc_datetime ASC
LIMIT $1 OFFSET $2;

-- name: GetCompletedMatches :many
SELECT 
    m.id, m.home_team_id, m.away_team_id, m.title, m.slug, m.tagline, m.utc_datetime, m.venue, m.home_score, m.away_score,
    ht.name AS home_team_name, ht.logo_url AS home_team_logo_url,
    at.name AS away_team_name, at.logo_url AS away_team_logo_url
FROM matches m
JOIN teams ht ON m.home_team_id = ht.id
JOIN teams at ON m.away_team_id = at.id
WHERE m.utc_datetime <= now()
ORDER BY m.utc_datetime DESC
LIMIT $1 OFFSET $2;

-- name: UpdateMatch :one
UPDATE matches
SET
    home_team_id = COALESCE($2, home_team_id),
    away_team_id = COALESCE($3, away_team_id),
    title = COALESCE($4, title),
    slug = COALESCE($5, slug),
    description = COALESCE($6, description),
    round = COALESCE($7, round),
    utc_datetime = COALESCE($8, utc_datetime),
    venue = COALESCE($9, venue),
    cover_image_url = COALESCE($10, cover_image_url),
    home_score = COALESCE($11, home_score),
    away_score = COALESCE($12, away_score),
    home_penalty_score = COALESCE($13, home_penalty_score),
    away_penalty_score = COALESCE($14, away_penalty_score),
    tagline = COALESCE($15, tagline),
    updated_at = now()
WHERE id = $1
RETURNING id, home_team_id, away_team_id, title, slug, description, round, utc_datetime, venue, cover_image_url, home_score, away_score, home_penalty_score, away_penalty_score, tagline, created_at, updated_at;

-- name: DeleteMatch :exec
DELETE FROM matches
WHERE id = $1;

-- name: GetMatchPerformances :many
SELECT 
    p.id, p.match_id, p.player_id, p.player_team_id, p.title, p.cover_image_url, p.jersey_number, p.is_starter, p.captain, p.minutes_played, p.average_rating, p.slug, p.tagline,
    pl.name AS player_name, pl.slug AS player_slug, pl.photo_url AS player_photo_url,
    t.name AS team_name, t.logo_url AS team_logo_url, t.flag_emoji,
    COALESCE((SELECT sv.value FROM performance_stats sv JOIN stat_types stt ON stt.id = sv.stat_type_id WHERE sv.performance_id = p.id AND stt.name = 'goals' LIMIT 1), 0)::int AS goals,
    COALESCE((SELECT sv.value FROM performance_stats sv JOIN stat_types stt ON stt.id = sv.stat_type_id WHERE sv.performance_id = p.id AND stt.name = 'assists' LIMIT 1), 0)::int AS assists,
    0::int AS dribbles,
    0::numeric(5,2) AS passes_accuracy,
    p.average_rating::numeric(3,1) AS average_rating,
    p.total_votes::bigint AS total_votes
FROM performances p
JOIN players pl ON p.player_id = pl.id
JOIN player_teams pt ON p.player_team_id = pt.id
JOIN teams t ON pt.team_id = t.id
WHERE p.match_id = $1
ORDER BY p.is_starter DESC, p.minutes_played DESC;

-- name: GetRelatedMatches :many
SELECT 
    m.id, m.home_team_id, m.away_team_id, m.title, m.slug, m.tagline, m.utc_datetime, m.home_score, m.away_score,
    ht.name AS home_team_name, ht.logo_url AS home_team_logo_url,
    at.name AS away_team_name, at.logo_url AS away_team_logo_url
FROM matches m
JOIN teams ht ON m.home_team_id = ht.id
JOIN teams at ON m.away_team_id = at.id
WHERE m.id <> $1
ORDER BY m.utc_datetime DESC
LIMIT $2;

-- name: GetOverviewStats :one
SELECT
    (SELECT COUNT(*) FROM matches)::bigint AS match_count,
    (SELECT COUNT(*) FROM performances)::bigint AS performance_count,
    (SELECT COUNT(*) FROM teams)::bigint AS team_count,
    ((SELECT COUNT(*) FROM match_ratings) + (SELECT COUNT(*) FROM performance_ratings))::bigint AS rating_count,
    ((SELECT COUNT(*) FROM match_reviews) + (SELECT COUNT(*) FROM performance_reviews))::bigint AS review_count;
