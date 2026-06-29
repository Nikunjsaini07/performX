-- name: CreateList :one
INSERT INTO lists (
    user_id,
    title,
    slug,
    description,
    cover_image_url,
    is_public
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING id, user_id, title, slug, description, cover_image_url, is_public, created_at, updated_at;

-- name: GetListByID :one
SELECT 
    l.id, l.user_id, l.title, l.slug, l.description, l.cover_image_url, l.is_public, l.created_at, l.updated_at,
    u.username AS author_username, u.display_name AS author_display_name, u.avatar_url AS author_avatar_url,
    COALESCE(likes.likes_count, 0)::bigint AS likes_count,
    COALESCE(items.items_count, 0)::bigint AS items_count
FROM lists l
JOIN users u ON l.user_id = u.id
LEFT JOIN (
    SELECT list_id, COUNT(*) AS likes_count
    FROM list_likes
    GROUP BY list_id
) likes ON l.id = likes.list_id
LEFT JOIN (
    SELECT list_id, COUNT(*) AS items_count
    FROM list_items
    GROUP BY list_id
) items ON l.id = items.list_id
WHERE l.id = $1 LIMIT 1;

-- name: GetListBySlug :one
SELECT 
    l.id, l.user_id, l.title, l.slug, l.description, l.cover_image_url, l.is_public, l.created_at, l.updated_at,
    u.username AS author_username, u.display_name AS author_display_name, u.avatar_url AS author_avatar_url,
    COALESCE(likes.likes_count, 0)::bigint AS likes_count,
    COALESCE(items.items_count, 0)::bigint AS items_count
FROM lists l
JOIN users u ON l.user_id = u.id
LEFT JOIN (
    SELECT list_id, COUNT(*) AS likes_count
    FROM list_likes
    GROUP BY list_id
) likes ON l.id = likes.list_id
LEFT JOIN (
    SELECT list_id, COUNT(*) AS items_count
    FROM list_items
    GROUP BY list_id
) items ON l.id = items.list_id
WHERE l.slug = $1 LIMIT 1;

-- name: GetUserLists :many
SELECT 
    l.id, l.user_id, l.title, l.slug, l.description, l.cover_image_url, l.is_public, l.created_at, l.updated_at,
    COALESCE(likes.likes_count, 0)::bigint AS likes_count,
    COALESCE(items.items_count, 0)::bigint AS items_count
FROM lists l
LEFT JOIN (
    SELECT list_id, COUNT(*) AS likes_count
    FROM list_likes
    GROUP BY list_id
) likes ON l.id = likes.list_id
LEFT JOIN (
    SELECT list_id, COUNT(*) AS items_count
    FROM list_items
    GROUP BY list_id
) items ON l.id = items.list_id
WHERE l.user_id = $1
ORDER BY l.created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetPublicLists :many
SELECT 
    l.id, l.user_id, l.title, l.slug, l.description, l.cover_image_url, l.is_public, l.created_at, l.updated_at,
    u.username AS author_username, u.display_name AS author_display_name, u.avatar_url AS author_avatar_url,
    COALESCE(likes.likes_count, 0)::bigint AS likes_count,
    COALESCE(items.items_count, 0)::bigint AS items_count
FROM lists l
JOIN users u ON l.user_id = u.id
LEFT JOIN (
    SELECT list_id, COUNT(*) AS likes_count
    FROM list_likes
    GROUP BY list_id
) likes ON l.id = likes.list_id
LEFT JOIN (
    SELECT list_id, COUNT(*) AS items_count
    FROM list_items
    GROUP BY list_id
) items ON l.id = items.list_id
WHERE l.is_public = TRUE
ORDER BY l.created_at DESC
LIMIT $1 OFFSET $2;

-- name: SearchLists :many
SELECT 
    l.id, l.user_id, l.title, l.slug, l.description, l.cover_image_url, l.is_public, l.created_at, l.updated_at,
    u.username AS author_username, u.display_name AS author_display_name, u.avatar_url AS author_avatar_url,
    COALESCE(likes.likes_count, 0)::bigint AS likes_count,
    COALESCE(items.items_count, 0)::bigint AS items_count
FROM lists l
JOIN users u ON l.user_id = u.id
LEFT JOIN (
    SELECT list_id, COUNT(*) AS likes_count
    FROM list_likes
    GROUP BY list_id
) likes ON l.id = likes.list_id
LEFT JOIN (
    SELECT list_id, COUNT(*) AS items_count
    FROM list_items
    GROUP BY list_id
) items ON l.id = items.list_id
WHERE l.is_public = TRUE AND (l.title ILIKE $1 OR l.description ILIKE $1)
ORDER BY l.created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateList :one
UPDATE lists
SET 
    title = $2,
    slug = $3,
    description = $4,
    updated_at = now()
WHERE id = $1
RETURNING id, user_id, title, slug, description, cover_image_url, is_public, created_at, updated_at;

-- name: UpdateListCoverImage :one
UPDATE lists
SET 
    cover_image_url = $2,
    updated_at = now()
WHERE id = $1
RETURNING id, user_id, title, slug, description, cover_image_url, is_public, created_at, updated_at;

-- name: UpdateListVisibility :one
UPDATE lists
SET 
    is_public = $2,
    updated_at = now()
WHERE id = $1
RETURNING id, user_id, title, slug, description, cover_image_url, is_public, created_at, updated_at;

-- name: DeleteList :exec
DELETE FROM lists
WHERE id = $1;

-- name: AddItemToList :one
INSERT INTO list_items (
    list_id,
    match_id,
    performance_id,
    position,
    note
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING id, list_id, match_id, performance_id, position, note, created_at;

-- name: RemoveItemFromList :exec
DELETE FROM list_items
WHERE id = $1;

-- name: UpdateListItemPosition :one
UPDATE list_items
SET position = $2
WHERE id = $1
RETURNING id, list_id, match_id, performance_id, position, note, created_at;

-- name: UpdateListItemNote :one
UPDATE list_items
SET note = $2
WHERE id = $1
RETURNING id, list_id, match_id, performance_id, position, note, created_at;

-- name: GetListItems :many
SELECT 
    li.id, li.list_id, li.match_id, li.performance_id, li.position, li.note, li.created_at,
    m.title AS match_title, m.slug AS match_slug, m.utc_datetime AS match_utc_datetime, m.home_score AS match_home_score, m.away_score AS match_away_score,
    ht.name AS home_team_name, ht.logo_url AS home_team_logo_url,
    at.name AS away_team_name, at.logo_url AS away_team_logo_url,
    p.title AS performance_title, p.provider_rating AS performance_provider_rating,
    pl.name AS player_name, pl.known_as AS player_known_as, pl.photo_url AS player_photo_url,
    pt.jersey_number AS player_jersey_number,
    t.name AS performance_team_name, t.logo_url AS performance_team_logo_url
FROM list_items li
-- Match joins
LEFT JOIN matches m ON li.match_id = m.id
LEFT JOIN teams ht ON m.home_team_id = ht.id
LEFT JOIN teams at ON m.away_team_id = at.id
-- Performance joins
LEFT JOIN performances p ON li.performance_id = p.id
LEFT JOIN players pl ON p.player_id = pl.id
LEFT JOIN player_teams pt ON p.player_team_id = pt.id
LEFT JOIN teams t ON pt.team_id = t.id
WHERE li.list_id = $1
ORDER BY li.position ASC;

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

-- name: GetTrendingLists :many
SELECT 
    l.id, l.user_id, l.title, l.slug, l.description, l.cover_image_url, l.is_public, l.created_at, l.updated_at,
    u.username AS author_username, u.display_name AS author_display_name, u.avatar_url AS author_avatar_url,
    COALESCE(likes.likes_count, 0)::bigint AS likes_count,
    COALESCE(items.items_count, 0)::bigint AS items_count
FROM lists l
JOIN users u ON l.user_id = u.id
LEFT JOIN (
    SELECT list_id, COUNT(*) AS likes_count
    FROM list_likes
    GROUP BY list_id
) likes ON l.id = likes.list_id
LEFT JOIN (
    SELECT list_id, COUNT(*) AS items_count
    FROM list_items
    GROUP BY list_id
) items ON l.id = items.list_id
WHERE l.is_public = TRUE
ORDER BY likes_count DESC, l.created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetRecentLists :many
SELECT 
    l.id, l.user_id, l.title, l.slug, l.description, l.cover_image_url, l.is_public, l.created_at, l.updated_at,
    u.username AS author_username, u.display_name AS author_display_name, u.avatar_url AS author_avatar_url,
    COALESCE(likes.likes_count, 0)::bigint AS likes_count,
    COALESCE(items.items_count, 0)::bigint AS items_count
FROM lists l
JOIN users u ON l.user_id = u.id
LEFT JOIN (
    SELECT list_id, COUNT(*) AS likes_count
    FROM list_likes
    GROUP BY list_id
) likes ON l.id = likes.list_id
LEFT JOIN (
    SELECT list_id, COUNT(*) AS items_count
    FROM list_items
    GROUP BY list_id
) items ON l.id = items.list_id
WHERE l.is_public = TRUE
ORDER BY l.created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetTopLikedLists :many
SELECT 
    l.id, l.user_id, l.title, l.slug, l.description, l.cover_image_url, l.is_public, l.created_at, l.updated_at,
    u.username AS author_username, u.display_name AS author_display_name, u.avatar_url AS author_avatar_url,
    COALESCE(likes.likes_count, 0)::bigint AS likes_count,
    COALESCE(items.items_count, 0)::bigint AS items_count
FROM lists l
JOIN users u ON l.user_id = u.id
LEFT JOIN (
    SELECT list_id, COUNT(*) AS likes_count
    FROM list_likes
    GROUP BY list_id
) likes ON l.id = likes.list_id
LEFT JOIN (
    SELECT list_id, COUNT(*) AS items_count
    FROM list_items
    GROUP BY list_id
) items ON l.id = items.list_id
WHERE l.is_public = TRUE
ORDER BY likes_count DESC
LIMIT $1 OFFSET $2;
