-- name: RatePerformance :one
INSERT INTO performance_ratings (
    performance_id,
    user_id,
    rating
) VALUES (
    $1, $2, $3
)
ON CONFLICT (performance_id, user_id) 
DO UPDATE SET 
    rating = EXCLUDED.rating,
    updated_at = now()
RETURNING id, performance_id, user_id, rating, created_at, updated_at;

-- name: UpdatePerformanceRating :one
UPDATE performance_ratings
SET 
    rating = $2,
    updated_at = now()
WHERE id = $1
RETURNING id, performance_id, user_id, rating, created_at, updated_at;

-- name: DeletePerformanceRating :exec
DELETE FROM performance_ratings
WHERE performance_id = $1 AND user_id = $2;

-- name: GetPerformanceRating :one
SELECT id, performance_id, user_id, rating, created_at, updated_at
FROM performance_ratings
WHERE id = $1 LIMIT 1;

-- name: GetUserPerformanceRating :one
SELECT id, performance_id, user_id, rating, created_at, updated_at
FROM performance_ratings
WHERE performance_id = $1 AND user_id = $2 LIMIT 1;

-- name: GetPerformanceAverageRating :one
-- Folds the provider seed rating in as the first vote so the live detail-page
-- value matches the stored average_rating column.
SELECT
    ((p.seed_rating + COALESCE(SUM(pr.rating), 0)) / (1 + COUNT(pr.id)))::numeric(3,1) AS average_rating,
    (1 + COUNT(pr.id))::bigint AS total_votes
FROM performances p
LEFT JOIN performance_ratings pr ON pr.performance_id = p.id
WHERE p.id = $1
GROUP BY p.seed_rating;

-- name: RefreshPerformanceRating :exec
-- Recomputes and persists the stored average_rating/total_votes for a
-- performance (seed vote + community votes). Call after any rating mutation.
UPDATE performances SET
    average_rating = sub.avg,
    total_votes = sub.cnt
FROM (
    SELECT ((p.seed_rating + COALESCE(SUM(pr.rating), 0)) / (1 + COUNT(pr.id)))::numeric(3,1) AS avg,
           (1 + COUNT(pr.id))::int AS cnt
    FROM performances p
    LEFT JOIN performance_ratings pr ON pr.performance_id = p.id
    WHERE p.id = $1
    GROUP BY p.seed_rating
) sub
WHERE performances.id = $1;

-- name: GetPerformanceRatingsCount :one
SELECT COUNT(*)::bigint AS ratings_count
FROM performance_ratings
WHERE performance_id = $1;

-- name: GetRecentlyRatedPerformances :many
SELECT 
    p.id, p.match_id, p.player_id, p.title, p.cover_image_url, p.average_rating,
    pl.name AS player_name, m.title AS match_title,
    latest_rating.max_created_at AS rated_at
FROM (
    SELECT performance_id, MAX(created_at) AS max_created_at
    FROM performance_ratings
    GROUP BY performance_id
) latest_rating
JOIN performances p ON latest_rating.performance_id = p.id
JOIN players pl ON p.player_id = pl.id
JOIN matches m ON p.match_id = m.id
ORDER BY rated_at DESC
LIMIT $1 OFFSET $2;

-- name: GetPerformanceRatings :many
SELECT 
    pr.id, pr.performance_id, pr.user_id, pr.rating, pr.created_at, pr.updated_at,
    u.username AS author_username, u.display_name AS author_display_name, u.avatar_url AS author_avatar_url
FROM performance_ratings pr
JOIN users u ON pr.user_id = u.id
WHERE pr.performance_id = $1
ORDER BY pr.created_at DESC
LIMIT $2 OFFSET $3;
