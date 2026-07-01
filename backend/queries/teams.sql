-- name: GetTeamByID :one
SELECT 
    t.id, t.sport_id, t.country_id, t.name, t.short_name, t.slug, t.type, t.logo_url, t.founded_year, t.created_at, t.updated_at,
    c.name AS country_name
FROM teams t
LEFT JOIN countries c ON t.country_id = c.id
WHERE t.id = $1 LIMIT 1;

-- name: GetTeamBySlug :one
SELECT 
    t.id, t.sport_id, t.country_id, t.name, t.short_name, t.slug, t.type, t.logo_url, t.founded_year, t.created_at, t.updated_at,
    c.name AS country_name
FROM teams t
LEFT JOIN countries c ON t.country_id = c.id
WHERE t.slug = $1 LIMIT 1;

-- name: SearchTeams :many
SELECT 
    t.id, t.sport_id, t.country_id, t.name, t.short_name, t.slug, t.type, t.logo_url
FROM teams t
WHERE t.name ILIKE $1 OR t.short_name ILIKE $1
ORDER BY t.name
LIMIT $2 OFFSET $3;

-- name: ListTeams :many
SELECT 
    t.id, t.sport_id, t.country_id, t.name, t.short_name, t.slug, t.type, t.logo_url, t.founded_year, t.created_at, t.updated_at
FROM teams t
ORDER BY t.name
LIMIT $1 OFFSET $2;

-- name: CreateTeam :one
INSERT INTO teams (
    sport_id,
    country_id,
    name,
    short_name,
    slug,
    type,
    logo_url,
    founded_year
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING id, sport_id, country_id, name, short_name, slug, type, logo_url, founded_year, created_at, updated_at;

-- name: UpdateTeam :one
UPDATE teams
SET
    country_id = COALESCE($2, country_id),
    name = COALESCE($3, name),
    short_name = COALESCE($4, short_name),
    slug = COALESCE($5, slug),
    type = COALESCE($6, type),
    logo_url = COALESCE($7, logo_url),
    founded_year = COALESCE($8, founded_year),
    updated_at = now()
WHERE id = $1
RETURNING id, sport_id, country_id, name, short_name, slug, type, logo_url, founded_year, created_at, updated_at;

-- name: DeleteTeam :exec
DELETE FROM teams
WHERE id = $1;

-- name: GetTeamsByCountry :many
SELECT 
    t.id, t.sport_id, t.country_id, t.name, t.short_name, t.slug, t.type, t.logo_url
FROM teams t
WHERE t.country_id = $1
ORDER BY t.name;



-- name: GetTeamMatches :many
SELECT 
    m.id, m.season_id, m.home_team_id, m.away_team_id, m.title, m.slug, m.utc_datetime, m.venue, m.home_score, m.away_score,
    ht.name AS home_team_name, ht.logo_url AS home_team_logo_url,
    at.name AS away_team_name, at.logo_url AS away_team_logo_url
FROM matches m
JOIN teams ht ON m.home_team_id = ht.id
JOIN teams at ON m.away_team_id = at.id
WHERE m.home_team_id = $1 OR m.away_team_id = $1
ORDER BY m.utc_datetime DESC
LIMIT $2 OFFSET $3;

-- name: GetTeamUpcomingMatches :many
SELECT 
    m.id, m.season_id, m.home_team_id, m.away_team_id, m.title, m.slug, m.utc_datetime, m.venue,
    ht.name AS home_team_name, ht.logo_url AS home_team_logo_url,
    at.name AS away_team_name, at.logo_url AS away_team_logo_url
FROM matches m
JOIN teams ht ON m.home_team_id = ht.id
JOIN teams at ON m.away_team_id = at.id
WHERE (m.home_team_id = $1 OR m.away_team_id = $1) AND m.utc_datetime > now()
ORDER BY m.utc_datetime ASC
LIMIT $2 OFFSET $3;

-- name: GetTeamCompletedMatches :many
SELECT 
    m.id, m.season_id, m.home_team_id, m.away_team_id, m.title, m.slug, m.utc_datetime, m.venue, m.home_score, m.away_score,
    ht.name AS home_team_name, ht.logo_url AS home_team_logo_url,
    at.name AS away_team_name, at.logo_url AS away_team_logo_url
FROM matches m
JOIN teams ht ON m.home_team_id = ht.id
JOIN teams at ON m.away_team_id = at.id
WHERE (m.home_team_id = $1 OR m.away_team_id = $1) AND m.utc_datetime <= now()
ORDER BY m.utc_datetime DESC
LIMIT $2 OFFSET $3;

-- name: GetTeamPerformances :many
SELECT 
    p.id AS performance_id, p.title AS performance_title, p.cover_image_url AS performance_cover_image, p.minutes_played, p.provider_rating,
    pl.id AS player_id, pl.name AS player_name, pl.slug AS player_slug, pl.photo_url AS player_photo_url,
    m.id AS match_id, m.title AS match_title, m.slug AS match_slug, m.utc_datetime AS match_utc_datetime
FROM performances p
JOIN player_teams pt ON p.player_team_id = pt.id
JOIN players pl ON p.player_id = pl.id
JOIN matches m ON p.match_id = m.id
WHERE pt.team_id = $1
ORDER BY m.utc_datetime DESC
LIMIT $2 OFFSET $3;

-- name: GetTeamAverageRating :one
SELECT 
    COALESCE(AVG(pr.rating), 0.0)::numeric(2,1) AS average_rating,
    COUNT(pr.id)::bigint AS total_votes
FROM performances p
JOIN player_teams pt ON p.player_team_id = pt.id
JOIN performance_ratings pr ON p.id = pr.performance_id
WHERE pt.team_id = $1;

-- name: GetTeamReviews :many
SELECT 
    r.id AS review_id, r.title AS review_title, r.content AS review_content, r.created_at AS review_created_at,
    u.id AS user_id, u.username, u.display_name, u.avatar_url,
    p.id AS performance_id, p.title AS performance_title,
    m.id AS match_id, m.title AS match_title
FROM performance_reviews r
JOIN performances p ON r.performance_id = p.id
JOIN player_teams pt ON p.player_team_id = pt.id
JOIN matches m ON p.match_id = m.id
JOIN users u ON r.user_id = u.id
WHERE pt.team_id = $1
ORDER BY r.created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetRelatedTeams :many
SELECT 
    t.id, t.name, t.short_name, t.slug, t.logo_url
FROM teams t
WHERE t.id <> $1 AND t.country_id = (SELECT country_id FROM teams WHERE id = $1) AND t.type = (SELECT type FROM teams WHERE id = $1)
ORDER BY t.name
LIMIT $2;

-- name: GetTeamAggregatedStats :many
SELECT 
    st.id AS stat_type_id,
    st.name AS stat_name,
    st.short_name AS stat_short_name,
    st.unit AS stat_unit,
    st.category AS stat_category,
    SUM(ms.value)::double precision AS total_value,
    AVG(ms.value)::double precision AS average_value
FROM match_stats ms
JOIN stat_types st ON ms.stat_type_id = st.id
WHERE ms.team_id = $1
GROUP BY st.id, st.name, st.short_name, st.unit, st.category, st.display_order
ORDER BY st.category, st.display_order;

-- name: GetTeamRatings :many
SELECT 
    mr.id AS rating_id, mr.rating, mr.created_at AS rated_at,
    u.id AS user_id, u.username, u.display_name, u.avatar_url,
    m.id AS match_id, m.title AS match_title, m.slug AS match_slug
FROM match_ratings mr
JOIN matches m ON mr.match_id = m.id
JOIN users u ON mr.user_id = u.id
WHERE m.home_team_id = $1 OR m.away_team_id = $1
ORDER BY mr.created_at DESC
LIMIT $2 OFFSET $3;
