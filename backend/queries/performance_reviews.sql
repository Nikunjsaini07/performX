-- name: CreatePerformanceReview :one
INSERT INTO performance_reviews (
    performance_id,
    user_id,
    title,
    content
) VALUES (
    $1, $2, $3, $4
)
RETURNING id, performance_id, user_id, title, content, created_at, updated_at;

-- name: GetPerformanceReviewByID :one
SELECT 
    pr.id, pr.performance_id, pr.user_id, pr.title, pr.content, pr.created_at, pr.updated_at,
    u.username AS author_username, u.display_name AS author_display_name, u.avatar_url AS author_avatar_url,
    p.title AS performance_title, pl.name AS player_name, m.title AS match_title,
    COALESCE(likes.likes_count, 0)::bigint AS likes_count
FROM performance_reviews pr
JOIN users u ON pr.user_id = u.id
JOIN performances p ON pr.performance_id = p.id
JOIN players pl ON p.player_id = pl.id
JOIN matches m ON p.match_id = m.id
LEFT JOIN (
    SELECT review_id, COUNT(*) AS likes_count
    FROM performance_review_likes
    GROUP BY review_id
) likes ON pr.id = likes.review_id
WHERE pr.id = $1 LIMIT 1;

-- name: GetPerformanceReviews :many
SELECT 
    pr.id, pr.performance_id, pr.user_id, pr.title, pr.content, pr.created_at, pr.updated_at,
    u.username AS author_username, u.display_name AS author_display_name, u.avatar_url AS author_avatar_url,
    COALESCE(likes.likes_count, 0)::bigint AS likes_count
FROM performance_reviews pr
JOIN users u ON pr.user_id = u.id
LEFT JOIN (
    SELECT review_id, COUNT(*) AS likes_count
    FROM performance_review_likes
    GROUP BY review_id
) likes ON pr.id = likes.review_id
WHERE pr.performance_id = $1
ORDER BY pr.created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetUserPerformanceReviews :many
SELECT 
    pr.id, pr.performance_id, pr.user_id, pr.title, pr.content, pr.created_at, pr.updated_at,
    p.title AS performance_title, pl.name AS player_name,
    COALESCE(likes.likes_count, 0)::bigint AS likes_count
FROM performance_reviews pr
JOIN performances p ON pr.performance_id = p.id
JOIN players pl ON p.player_id = pl.id
LEFT JOIN (
    SELECT review_id, COUNT(*) AS likes_count
    FROM performance_review_likes
    GROUP BY review_id
) likes ON pr.id = likes.review_id
WHERE pr.user_id = $1
ORDER BY pr.created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdatePerformanceReview :one
UPDATE performance_reviews
SET 
    title = COALESCE($2, title),
    content = COALESCE($3, content),
    updated_at = now()
WHERE id = $1
RETURNING id, performance_id, user_id, title, content, created_at, updated_at;

-- name: DeletePerformanceReview :exec
DELETE FROM performance_reviews
WHERE id = $1;

-- name: GetPerformanceReviewCount :one
SELECT COUNT(*)::bigint AS review_count
FROM performance_reviews
WHERE performance_id = $1;

-- name: GetPerformanceReviewAverageRating :one
SELECT COALESCE(AVG(r.rating), 0.0)::numeric(2,1) AS average_rating
FROM performance_reviews pr
JOIN performance_ratings r ON pr.performance_id = r.performance_id AND pr.user_id = r.user_id
WHERE pr.performance_id = $1;

-- name: GetRecentPerformanceReviews :many
SELECT 
    pr.id, pr.performance_id, pr.user_id, pr.title, pr.content, pr.created_at, pr.updated_at,
    u.username AS author_username, u.display_name AS author_display_name, u.avatar_url AS author_avatar_url,
    p.title AS performance_title, pl.name AS player_name,
    COALESCE(likes.likes_count, 0)::bigint AS likes_count
FROM performance_reviews pr
JOIN users u ON pr.user_id = u.id
JOIN performances p ON pr.performance_id = p.id
JOIN players pl ON p.player_id = pl.id
LEFT JOIN (
    SELECT review_id, COUNT(*) AS likes_count
    FROM performance_review_likes
    GROUP BY review_id
) likes ON pr.id = likes.review_id
ORDER BY pr.created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetPopularPerformanceReviews :many
SELECT 
    pr.id, pr.performance_id, pr.user_id, pr.title, pr.content, pr.created_at, pr.updated_at,
    u.username AS author_username, u.display_name AS author_display_name, u.avatar_url AS author_avatar_url,
    p.title AS performance_title, pl.name AS player_name,
    COALESCE(likes.likes_count, 0)::bigint AS likes_count
FROM performance_reviews pr
JOIN users u ON pr.user_id = u.id
JOIN performances p ON pr.performance_id = p.id
JOIN players pl ON p.player_id = pl.id
LEFT JOIN (
    SELECT review_id, COUNT(*) AS likes_count
    FROM performance_review_likes
    GROUP BY review_id
) likes ON pr.id = likes.review_id
ORDER BY likes_count DESC, pr.created_at DESC
LIMIT $1 OFFSET $2;

-- name: HasUserReviewedPerformance :one
SELECT EXISTS (
    SELECT 1 FROM performance_reviews
    WHERE performance_id = $1 AND user_id = $2
)::boolean AS has_reviewed;
