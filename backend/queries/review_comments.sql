-- ============================================================================
-- Match Review Comments
-- ============================================================================

-- name: CreateMatchReviewComment :one
INSERT INTO match_review_comments (
    review_id,
    user_id,
    body
) VALUES (
    $1, $2, $3
)
RETURNING id, review_id, user_id, body, created_at, updated_at;

-- name: GetMatchReviewCommentByID :one
SELECT 
    c.id, c.review_id, c.user_id, c.body, c.created_at, c.updated_at,
    u.username, u.display_name, u.avatar_url
FROM match_review_comments c
JOIN users u ON c.user_id = u.id
WHERE c.id = $1 LIMIT 1;

-- name: GetMatchReviewComments :many
SELECT 
    c.id, c.review_id, c.user_id, c.body, c.created_at, c.updated_at,
    u.username, u.display_name, u.avatar_url
FROM match_review_comments c
JOIN users u ON c.user_id = u.id
WHERE c.review_id = $1
ORDER BY c.created_at ASC
LIMIT $2 OFFSET $3;

-- name: GetUserMatchReviewComments :many
SELECT 
    c.id, c.review_id, c.user_id, c.body, c.created_at, c.updated_at,
    r.match_id, m.title AS match_title
FROM match_review_comments c
JOIN match_reviews r ON c.review_id = r.id
JOIN matches m ON r.match_id = m.id
WHERE c.user_id = $1
ORDER BY c.created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateMatchReviewComment :one
UPDATE match_review_comments
SET 
    body = $2,
    updated_at = now()
WHERE id = $1
RETURNING id, review_id, user_id, body, created_at, updated_at;

-- name: DeleteMatchReviewComment :exec
DELETE FROM match_review_comments
WHERE id = $1;

-- name: GetMatchReviewCommentCount :one
SELECT COUNT(*)::bigint AS comment_count
FROM match_review_comments
WHERE review_id = $1;

-- name: GetRecentMatchReviewComments :many
SELECT 
    c.id, c.review_id, c.user_id, c.body, c.created_at, c.updated_at,
    u.username, u.display_name, u.avatar_url,
    m.title AS match_title
FROM match_review_comments c
JOIN users u ON c.user_id = u.id
JOIN match_reviews r ON c.review_id = r.id
JOIN matches m ON r.match_id = m.id
ORDER BY c.created_at DESC
LIMIT $1 OFFSET $2;


-- ============================================================================
-- Performance Review Comments
-- ============================================================================

-- name: CreatePerformanceReviewComment :one
INSERT INTO performance_review_comments (
    review_id,
    user_id,
    body
) VALUES (
    $1, $2, $3
)
RETURNING id, review_id, user_id, body, created_at, updated_at;

-- name: GetPerformanceReviewCommentByID :one
SELECT 
    c.id, c.review_id, c.user_id, c.body, c.created_at, c.updated_at,
    u.username, u.display_name, u.avatar_url
FROM performance_review_comments c
JOIN users u ON c.user_id = u.id
WHERE c.id = $1 LIMIT 1;

-- name: GetPerformanceReviewComments :many
SELECT 
    c.id, c.review_id, c.user_id, c.body, c.created_at, c.updated_at,
    u.username, u.display_name, u.avatar_url
FROM performance_review_comments c
JOIN users u ON c.user_id = u.id
WHERE c.review_id = $1
ORDER BY c.created_at ASC
LIMIT $2 OFFSET $3;

-- name: GetUserPerformanceReviewComments :many
SELECT 
    c.id, c.review_id, c.user_id, c.body, c.created_at, c.updated_at,
    r.performance_id, p.title AS performance_title
FROM performance_review_comments c
JOIN performance_reviews r ON c.review_id = r.id
JOIN performances p ON r.performance_id = p.id
WHERE c.user_id = $1
ORDER BY c.created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdatePerformanceReviewComment :one
UPDATE performance_review_comments
SET 
    body = $2,
    updated_at = now()
WHERE id = $1
RETURNING id, review_id, user_id, body, created_at, updated_at;

-- name: DeletePerformanceReviewComment :exec
DELETE FROM performance_review_comments
WHERE id = $1;

-- name: GetPerformanceReviewCommentCount :one
SELECT COUNT(*)::bigint AS comment_count
FROM performance_review_comments
WHERE review_id = $1;

-- name: GetRecentPerformanceReviewComments :many
SELECT 
    c.id, c.review_id, c.user_id, c.body, c.created_at, c.updated_at,
    u.username, u.display_name, u.avatar_url,
    p.title AS performance_title
FROM performance_review_comments c
JOIN users u ON c.user_id = u.id
JOIN performance_reviews r ON c.review_id = r.id
JOIN performances p ON r.performance_id = p.id
ORDER BY c.created_at DESC
LIMIT $1 OFFSET $2;
