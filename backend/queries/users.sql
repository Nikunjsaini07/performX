-- name: CreateUser :one
INSERT INTO users (
    username,
    display_name,
    email,
    bio,
    avatar_url
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING id, username, display_name, email, bio, avatar_url, created_at, updated_at;

-- name: GetUserByID :one
SELECT id, username, display_name, email, bio, avatar_url, created_at, updated_at
FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByUsername :one
SELECT id, username, display_name, email, bio, avatar_url, created_at, updated_at
FROM users
WHERE username = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT id, username, display_name, email, bio, avatar_url, created_at, updated_at
FROM users
WHERE email = $1 LIMIT 1;

-- name: SearchUsers :many
SELECT id, username, display_name, email, bio, avatar_url, created_at, updated_at
FROM users
WHERE username ILIKE $1 OR display_name ILIKE $1
ORDER BY username
LIMIT $2 OFFSET $3;

-- name: UpdateUsername :one
UPDATE users
SET username = $2, updated_at = now()
WHERE id = $1
RETURNING id, username, display_name, email, bio, avatar_url, created_at, updated_at;

-- name: UpdateDisplayName :one
UPDATE users
SET display_name = $2, updated_at = now()
WHERE id = $1
RETURNING id, username, display_name, email, bio, avatar_url, created_at, updated_at;

-- name: UpdateBio :one
UPDATE users
SET bio = $2, updated_at = now()
WHERE id = $1
RETURNING id, username, display_name, email, bio, avatar_url, created_at, updated_at;

-- name: UpdateAvatar :one
UPDATE users
SET avatar_url = $2, updated_at = now()
WHERE id = $1
RETURNING id, username, display_name, email, bio, avatar_url, created_at, updated_at;

-- name: UpdateEmail :one
UPDATE users
SET email = $2, updated_at = now()
WHERE id = $1
RETURNING id, username, display_name, email, bio, avatar_url, created_at, updated_at;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: GetPublicProfile :one
SELECT id, username, display_name, bio, avatar_url, created_at
FROM users
WHERE id = $1 LIMIT 1;

-- name: GetFollowers :many
SELECT 
    u.id, u.username, u.display_name, u.email, u.bio, u.avatar_url, uf.created_at AS followed_at
FROM user_follows uf
JOIN users u ON uf.follower_id = u.id
WHERE uf.following_id = $1
ORDER BY uf.created_at DESC;

-- name: GetFollowing :many
SELECT 
    u.id, u.username, u.display_name, u.email, u.bio, u.avatar_url, uf.created_at AS followed_at
FROM user_follows uf
JOIN users u ON uf.following_id = u.id
WHERE uf.follower_id = $1
ORDER BY uf.created_at DESC;

-- name: IsFollowing :one
SELECT EXISTS(
    SELECT 1 FROM user_follows
    WHERE follower_id = $1 AND following_id = $2
)::boolean AS is_following;

-- name: GetFollowersCount :one
SELECT COUNT(*)::bigint AS followers_count
FROM user_follows
WHERE following_id = $1;

-- name: GetFollowingCount :one
SELECT COUNT(*)::bigint AS following_count
FROM user_follows
WHERE follower_id = $1;

-- name: GetUserReviewCount :one
SELECT (
    (SELECT COUNT(*) FROM match_reviews mr WHERE mr.user_id = $1) +
    (SELECT COUNT(*) FROM performance_reviews pr WHERE pr.user_id = $1)
)::bigint AS review_count;

-- name: GetUserRatingCount :one
SELECT (
    (SELECT COUNT(*) FROM match_ratings mr WHERE mr.user_id = $1) +
    (SELECT COUNT(*) FROM performance_ratings pr WHERE pr.user_id = $1)
)::bigint AS rating_count;

-- name: GetUserListCount :one
SELECT COUNT(*)::bigint AS list_count
FROM lists l
WHERE l.user_id = $1;

-- name: GetUserLikesReceived :one
SELECT (
    (SELECT COUNT(*) FROM list_likes ll JOIN lists l ON ll.list_id = l.id WHERE l.user_id = $1) +
    (SELECT COUNT(*) FROM match_review_likes mrl JOIN match_reviews mr ON mrl.review_id = mr.id WHERE mr.user_id = $1) +
    (SELECT COUNT(*) FROM performance_review_likes prl JOIN performance_reviews pr ON prl.review_id = pr.id WHERE pr.user_id = $1) +
    (SELECT COUNT(*) FROM match_review_comment_likes mrcl JOIN match_review_comments mrc ON mrcl.comment_id = mrc.id WHERE mrc.user_id = $1)
)::bigint AS likes_count;

-- name: GetUserRecentReviews :many
SELECT 'match' AS review_type, mr.id, mr.match_id AS entity_id, mr.title, mr.content, mr.created_at, mr.updated_at
FROM match_reviews mr
WHERE mr.user_id = $1
UNION ALL
SELECT 'performance' AS review_type, pr.id, pr.performance_id AS entity_id, pr.title, pr.content, pr.created_at, pr.updated_at
FROM performance_reviews pr
WHERE pr.user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetUserRecentRatings :many
SELECT 'match' AS rating_type, mr.id, mr.match_id AS entity_id, mr.rating, mr.created_at, mr.updated_at
FROM match_ratings mr
WHERE mr.user_id = $1
UNION ALL
SELECT 'performance' AS rating_type, pr.id, pr.performance_id AS entity_id, pr.rating, pr.created_at, pr.updated_at
FROM performance_ratings pr
WHERE pr.user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetUserRecentComments :many
SELECT 'match_review' AS comment_type, mrc.id, mrc.review_id, mrc.body, mrc.created_at, mrc.updated_at
FROM match_review_comments mrc
WHERE mrc.user_id = $1
UNION ALL
SELECT 'performance_review' AS comment_type, prc.id, prc.review_id, prc.body, prc.created_at, prc.updated_at
FROM performance_review_comments prc
WHERE prc.user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetUserRecentLists :many
SELECT l.id, l.user_id, l.title, l.description, l.cover_image_url, l.is_public, l.created_at, l.updated_at
FROM lists l
WHERE l.user_id = $1
ORDER BY l.created_at DESC
LIMIT $2 OFFSET $3;
