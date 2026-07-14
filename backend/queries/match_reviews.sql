-- name: CreateMatchReview :one
INSERT INTO match_reviews (
    match_id,
    user_id,
    title,
    content
) VALUES (
    $1, $2, $3, $4
)
RETURNING id, match_id, user_id, title, content, created_at, updated_at;

-- name: GetMatchReviewByID :one
SELECT 
    mr.id, mr.match_id, mr.user_id, mr.title, mr.content, mr.created_at, mr.updated_at,
    u.username AS author_username, u.display_name AS author_display_name, u.avatar_url AS author_avatar_url,
    m.title AS match_title, m.slug AS match_slug,
    COALESCE(likes.likes_count, 0)::bigint AS likes_count,
    COALESCE(comments.comment_count, 0)::bigint AS comment_count
FROM match_reviews mr
JOIN users u ON mr.user_id = u.id
JOIN matches m ON mr.match_id = m.id
LEFT JOIN (
    SELECT review_id, COUNT(*) AS likes_count
    FROM match_review_likes
    GROUP BY review_id
) likes ON mr.id = likes.review_id
LEFT JOIN (
    SELECT review_id, COUNT(*) AS comment_count
    FROM match_review_comments
    GROUP BY review_id
) comments ON mr.id = comments.review_id
WHERE mr.id = $1 LIMIT 1;

-- name: GetMatchReviews :many
SELECT 
    mr.id, mr.match_id, mr.user_id, mr.title, mr.content, mr.created_at, mr.updated_at,
    u.username AS author_username, u.display_name AS author_display_name, u.avatar_url AS author_avatar_url,
    COALESCE(likes.likes_count, 0)::bigint AS likes_count,
    COALESCE(comments.comment_count, 0)::bigint AS comment_count,
    CASE WHEN $4::uuid IS NOT NULL THEN 
        EXISTS(SELECT 1 FROM match_review_likes WHERE review_id = mr.id AND user_id = $4)
    ELSE false END AS liked_by_me
FROM match_reviews mr
JOIN users u ON mr.user_id = u.id
LEFT JOIN (
    SELECT review_id, COUNT(*) AS likes_count
    FROM match_review_likes
    GROUP BY review_id
) likes ON mr.id = likes.review_id
LEFT JOIN (
    SELECT review_id, COUNT(*) AS comment_count
    FROM match_review_comments
    GROUP BY review_id
) comments ON mr.id = comments.review_id
WHERE mr.match_id = $1
ORDER BY mr.created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetUserMatchReviews :many
SELECT 
    mr.id, mr.match_id, mr.user_id, mr.title, mr.content, mr.created_at, mr.updated_at,
    m.title AS match_title, m.slug AS match_slug,
    COALESCE(likes.likes_count, 0)::bigint AS likes_count
FROM match_reviews mr
JOIN matches m ON mr.match_id = m.id
LEFT JOIN (
    SELECT review_id, COUNT(*) AS likes_count
    FROM match_review_likes
    GROUP BY review_id
) likes ON mr.id = likes.review_id
WHERE mr.user_id = $1
ORDER BY mr.created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateMatchReview :one
UPDATE match_reviews
SET 
    title = COALESCE($2, title),
    content = COALESCE($3, content),
    updated_at = now()
WHERE id = $1
RETURNING id, match_id, user_id, title, content, created_at, updated_at;

-- name: DeleteMatchReview :exec
DELETE FROM match_reviews
WHERE id = $1;

-- name: GetMatchReviewCount :one
SELECT COUNT(*)::bigint AS review_count
FROM match_reviews
WHERE match_id = $1;

-- name: GetMatchReviewAverageRating :one
SELECT COALESCE(AVG(r.rating), 0.0)::numeric(2,1) AS average_rating
FROM match_reviews mr
JOIN match_ratings r ON mr.match_id = r.match_id AND mr.user_id = r.user_id
WHERE mr.match_id = $1;

-- name: GetRecentMatchReviews :many
SELECT 
    mr.id, mr.match_id, mr.user_id, mr.title, mr.content, mr.created_at, mr.updated_at,
    u.username AS author_username, u.display_name AS author_display_name, u.avatar_url AS author_avatar_url,
    m.title AS match_title, m.slug AS match_slug,
    COALESCE(likes.likes_count, 0)::bigint AS likes_count,
    COALESCE(comments.comment_count, 0)::bigint AS comment_count
FROM match_reviews mr
JOIN users u ON mr.user_id = u.id
JOIN matches m ON mr.match_id = m.id
LEFT JOIN (
    SELECT review_id, COUNT(*) AS likes_count
    FROM match_review_likes
    GROUP BY review_id
) likes ON mr.id = likes.review_id
LEFT JOIN (
    SELECT review_id, COUNT(*) AS comment_count
    FROM match_review_comments
    GROUP BY review_id
) comments ON mr.id = comments.review_id
ORDER BY mr.created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetPopularMatchReviews :many
SELECT 
    mr.id, mr.match_id, mr.user_id, mr.title, mr.content, mr.created_at, mr.updated_at,
    u.username AS author_username, u.display_name AS author_display_name, u.avatar_url AS author_avatar_url,
    m.title AS match_title, m.slug AS match_slug,
    COALESCE(likes.likes_count, 0)::bigint AS likes_count,
    COALESCE(comments.comment_count, 0)::bigint AS comment_count
FROM match_reviews mr
JOIN users u ON mr.user_id = u.id
JOIN matches m ON mr.match_id = m.id
LEFT JOIN (
    SELECT review_id, COUNT(*) AS likes_count
    FROM match_review_likes
    GROUP BY review_id
) likes ON mr.id = likes.review_id
LEFT JOIN (
    SELECT review_id, COUNT(*) AS comment_count
    FROM match_review_comments
    GROUP BY review_id
) comments ON mr.id = comments.review_id
ORDER BY likes_count DESC, mr.created_at DESC
LIMIT $1 OFFSET $2;

-- name: HasUserReviewedMatch :one
SELECT EXISTS (
    SELECT 1 FROM match_reviews
    WHERE match_id = $1 AND user_id = $2
)::boolean AS has_reviewed;
