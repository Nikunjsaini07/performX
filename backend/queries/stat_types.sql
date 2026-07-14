-- name: CreateStatType :one
INSERT INTO stat_types (
    name,
    short_name,
    unit,
    category,
    display_order
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING id, name, short_name, unit, category, display_order;

-- name: GetStatTypeByID :one
SELECT id, name, short_name, unit, category, display_order
FROM stat_types
WHERE id = $1 LIMIT 1;

-- name: ListStatTypes :many
SELECT id, name, short_name, unit, category, display_order
FROM stat_types
ORDER BY category, display_order;

-- name: GetStatTypesByCategory :many
SELECT id, name, short_name, unit, category, display_order
FROM stat_types
WHERE category = $1
ORDER BY display_order;

-- name: UpdateStatType :one
UPDATE stat_types
SET 
    name = COALESCE($2, name),
    short_name = COALESCE($3, short_name),
    unit = COALESCE($4, unit),
    category = COALESCE($5, category),
    display_order = COALESCE($6, display_order)
WHERE id = $1
RETURNING id, name, short_name, unit, category, display_order;

-- name: DeleteStatType :exec
DELETE FROM stat_types
WHERE id = $1;
