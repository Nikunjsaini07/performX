-- name: FollowUser :one
INSERT INTO user_follows (
    follower_id,
    following_id
) VALUES (
    $1, $2
)
RETURNING follower_id, following_id, created_at;

-- name: UnfollowUser :exec
DELETE FROM user_follows
WHERE follower_id = $1 AND following_id = $2;

-- name: IsFollowing :one
SELECT EXISTS(
    SELECT 1 FROM user_follows
    WHERE follower_id = $1 AND following_id = $2
)::boolean AS is_following;

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

-- name: GetFollowersCount :one
SELECT COUNT(*)::bigint AS followers_count
FROM user_follows
WHERE following_id = $1;

-- name: GetFollowingCount :one
SELECT COUNT(*)::bigint AS following_count
FROM user_follows
WHERE follower_id = $1;

-- name: GetMutualFollowers :many
SELECT 
    u.id, u.username, u.display_name, u.avatar_url
FROM users u
WHERE u.id IN (
    SELECT uf1.follower_id FROM user_follows uf1 WHERE uf1.following_id = $1
    INTERSECT
    SELECT uf2.follower_id FROM user_follows uf2 WHERE uf2.following_id = $2
)
ORDER BY u.username;

-- name: GetMutualFollowing :many
SELECT 
    u.id, u.username, u.display_name, u.avatar_url
FROM users u
WHERE u.id IN (
    SELECT uf1.following_id FROM user_follows uf1 WHERE uf1.follower_id = $1
    INTERSECT
    SELECT uf2.following_id FROM user_follows uf2 WHERE uf2.follower_id = $2
)
ORDER BY u.username;
