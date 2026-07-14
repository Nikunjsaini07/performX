-- name: RateMatch :one
INSERT INTO match_ratings (
    match_id,
    user_id,
    rating
) VALUES (
    $1, $2, $3
)
ON CONFLICT (match_id, user_id) 
DO UPDATE SET 
    rating = EXCLUDED.rating,
    updated_at = now()
RETURNING id, match_id, user_id, rating, created_at, updated_at;

-- name: UpdateMatchRating :one
UPDATE match_ratings
SET 
    rating = $2,
    updated_at = now()
WHERE id = $1
RETURNING id, match_id, user_id, rating, created_at, updated_at;

-- name: DeleteMatchRating :exec
DELETE FROM match_ratings
WHERE match_id = $1 AND user_id = $2;

-- name: GetMatchRating :one
SELECT id, match_id, user_id, rating, created_at, updated_at
FROM match_ratings
WHERE id = $1 LIMIT 1;

-- name: GetUserMatchRating :one
SELECT id, match_id, user_id, rating, created_at, updated_at
FROM match_ratings
WHERE match_id = $1 AND user_id = $2 LIMIT 1;

-- name: GetMatchAverageRating :one
SELECT 
    COALESCE(AVG(rating), 0.0)::numeric(3,1) AS average_rating,
    COUNT(id)::bigint AS total_votes
FROM match_ratings
WHERE match_id = $1;

-- name: RefreshMatchRating :exec
-- Recomputes and persists the stored average_rating/total_votes for a match.
-- Call after any rating mutation. Matches have no provider seed.
UPDATE matches SET
    average_rating = COALESCE((SELECT AVG(rating) FROM match_ratings mr WHERE mr.match_id = $1), 0)::numeric(3,1),
    total_votes = (SELECT COUNT(*) FROM match_ratings mr WHERE mr.match_id = $1)::int
WHERE id = $1;

-- name: GetMatchRatingsCount :one
SELECT COUNT(*)::bigint AS ratings_count
FROM match_ratings
WHERE match_id = $1;

-- name: GetTopRatedMatches :many
SELECT 
    m.id, m.home_team_id, m.away_team_id, m.title, m.slug, m.utc_datetime, m.home_score, m.away_score,
    ht.name AS home_team_name, ht.logo_url AS home_team_logo_url,
    at.name AS away_team_name, at.logo_url AS away_team_logo_url,
    m.average_rating,
    m.total_votes
FROM matches m
JOIN teams ht ON m.home_team_id = ht.id
JOIN teams at ON m.away_team_id = at.id
ORDER BY m.average_rating DESC, m.utc_datetime DESC
LIMIT $1 OFFSET $2;

-- name: GetRecentlyRatedMatches :many
SELECT 
    m.id, m.home_team_id, m.away_team_id, m.title, m.slug, m.utc_datetime, m.home_score, m.away_score,
    ht.name AS home_team_name, ht.logo_url AS home_team_logo_url,
    at.name AS away_team_name, at.logo_url AS away_team_logo_url,
    latest_rating.max_created_at AS rated_at
FROM (
    SELECT match_id, MAX(created_at) AS max_created_at
    FROM match_ratings
    GROUP BY match_id
) latest_rating
JOIN matches m ON latest_rating.match_id = m.id
JOIN teams ht ON m.home_team_id = ht.id
JOIN teams at ON m.away_team_id = at.id
ORDER BY rated_at DESC
LIMIT $1 OFFSET $2;

-- name: GetMatchRatings :many
SELECT 
    mr.id, mr.match_id, mr.user_id, mr.rating, mr.created_at, mr.updated_at,
    u.username AS author_username, u.display_name AS author_display_name, u.avatar_url AS author_avatar_url
FROM match_ratings mr
JOIN users u ON mr.user_id = u.id
WHERE mr.match_id = $1
ORDER BY mr.created_at DESC
LIMIT $2 OFFSET $3;
