-- name: LikeList :one
INSERT INTO list_likes (
    list_id,
    user_id
) VALUES (
    $1, $2
)
RETURNING list_id, user_id, created_at;

-- name: UnlikeList :exec
DELETE FROM list_likes
WHERE list_id = $1 AND user_id = $2;

-- name: HasUserLikedList :one
SELECT EXISTS (
    SELECT 1 FROM list_likes
    WHERE list_id = $1 AND user_id = $2
)::boolean AS has_liked;

-- name: GetListLikesCount :one
SELECT COUNT(*)::bigint AS likes_count
FROM list_likes
WHERE list_id = $1;
