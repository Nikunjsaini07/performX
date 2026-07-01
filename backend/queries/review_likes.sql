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


-- ============================================================================
-- Match Review Comment Likes
-- ============================================================================

-- name: LikeMatchReviewComment :one
INSERT INTO match_review_comment_likes (
    comment_id,
    user_id
) VALUES (
    $1, $2
)
RETURNING comment_id, user_id, created_at;

-- name: UnlikeMatchReviewComment :exec
DELETE FROM match_review_comment_likes
WHERE comment_id = $1 AND user_id = $2;

-- name: HasUserLikedMatchReviewComment :one
SELECT EXISTS (
    SELECT 1 FROM match_review_comment_likes
    WHERE comment_id = $1 AND user_id = $2
)::boolean AS has_liked;

-- name: GetMatchReviewCommentLikesCount :one
SELECT COUNT(*)::bigint AS likes_count
FROM match_review_comment_likes
WHERE comment_id = $1;

-- name: GetUsersWhoLikedMatchReviewComment :many
SELECT 
    u.id, u.username, u.display_name, u.avatar_url, mrcl.created_at AS liked_at
FROM match_review_comment_likes mrcl
JOIN users u ON mrcl.user_id = u.id
WHERE mrcl.comment_id = $1
ORDER BY mrcl.created_at DESC;


-- ============================================================================
-- Performance Review Comment Likes
-- ============================================================================

-- name: LikePerformanceReviewComment :one
INSERT INTO performance_review_comment_likes (
    comment_id,
    user_id
) VALUES (
    $1, $2
)
RETURNING comment_id, user_id, created_at;

-- name: UnlikePerformanceReviewComment :exec
DELETE FROM performance_review_comment_likes
WHERE comment_id = $1 AND user_id = $2;

-- name: HasUserLikedPerformanceReviewComment :one
SELECT EXISTS (
    SELECT 1 FROM performance_review_comment_likes
    WHERE comment_id = $1 AND user_id = $2
)::boolean AS has_liked;

-- name: GetPerformanceReviewCommentLikesCount :one
SELECT COUNT(*)::bigint AS likes_count
FROM performance_review_comment_likes
WHERE comment_id = $1;

-- name: GetUsersWhoLikedPerformanceReviewComment :many
SELECT 
    u.id, u.username, u.display_name, u.avatar_url, prcl.created_at AS liked_at
FROM performance_review_comment_likes prcl
JOIN users u ON prcl.user_id = u.id
WHERE prcl.comment_id = $1
ORDER BY prcl.created_at DESC;
