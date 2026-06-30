-- ============================================================================
-- Match Review Likes
-- ============================================================================

-- name: LikeMatchReview :one
INSERT INTO match_review_likes (
    review_id,
    user_id
) VALUES (
    $1, $2
)
RETURNING review_id, user_id, created_at;

-- name: UnlikeMatchReview :exec
DELETE FROM match_review_likes
WHERE review_id = $1 AND user_id = $2;

-- name: HasUserLikedMatchReview :one
SELECT EXISTS (
    SELECT 1 FROM match_review_likes
    WHERE review_id = $1 AND user_id = $2
)::boolean AS has_liked;

-- name: GetMatchReviewLikesCount :one
SELECT COUNT(*)::bigint AS likes_count
FROM match_review_likes
WHERE review_id = $1;

-- name: GetUsersWhoLikedMatchReview :many
SELECT 
    u.id, u.username, u.display_name, u.avatar_url, mrl.created_at AS liked_at
FROM match_review_likes mrl
JOIN users u ON mrl.user_id = u.id
WHERE mrl.review_id = $1
ORDER BY mrl.created_at DESC;


-- ============================================================================
-- Performance Review Likes
-- ============================================================================

-- name: LikePerformanceReview :one
INSERT INTO performance_review_likes (
    review_id,
    user_id
) VALUES (
    $1, $2
)
RETURNING review_id, user_id, created_at;

-- name: UnlikePerformanceReview :exec
DELETE FROM performance_review_likes
WHERE review_id = $1 AND user_id = $2;

-- name: HasUserLikedPerformanceReview :one
SELECT EXISTS (
    SELECT 1 FROM performance_review_likes
    WHERE review_id = $1 AND user_id = $2
)::boolean AS has_liked;

-- name: GetPerformanceReviewLikesCount :one
SELECT COUNT(*)::bigint AS likes_count
FROM performance_review_likes
WHERE review_id = $1;

-- name: GetUsersWhoLikedPerformanceReview :many
SELECT 
    u.id, u.username, u.display_name, u.avatar_url, prl.created_at AS liked_at
FROM performance_review_likes prl
JOIN users u ON prl.user_id = u.id
WHERE prl.review_id = $1
ORDER BY prl.created_at DESC;
