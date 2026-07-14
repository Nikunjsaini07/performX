-- name: GetPerformanceStats :many
SELECT 
    ps.id, ps.performance_id, ps.stat_type_id, ps.value,
    st.name AS stat_name, st.short_name AS stat_short_name, st.unit AS stat_unit, st.category AS stat_category, st.display_order AS stat_display_order
FROM performance_stats ps
JOIN stat_types st ON ps.stat_type_id = st.id
WHERE ps.performance_id = $1
ORDER BY st.category, st.display_order;

-- name: GetPerformanceStat :one
SELECT 
    ps.id, ps.performance_id, ps.stat_type_id, ps.value,
    st.name AS stat_name, st.short_name AS stat_short_name, st.unit AS stat_unit, st.category AS stat_category
FROM performance_stats ps
JOIN stat_types st ON ps.stat_type_id = st.id
WHERE ps.performance_id = $1 AND ps.stat_type_id = $2 LIMIT 1;

-- name: CreatePerformanceStat :one
INSERT INTO performance_stats (
    performance_id,
    stat_type_id,
    value
) VALUES (
    $1, $2, $3
)
RETURNING id, performance_id, stat_type_id, value, created_at;

-- name: CreatePerformanceStats :copyfrom
INSERT INTO performance_stats (
    performance_id,
    stat_type_id,
    value
) VALUES (
    $1, $2, $3
);

-- name: UpdatePerformanceStat :one
UPDATE performance_stats
SET value = $3
WHERE performance_id = $1 AND stat_type_id = $2
RETURNING id, performance_id, stat_type_id, value, created_at;

-- name: DeletePerformanceStat :exec
DELETE FROM performance_stats
WHERE performance_id = $1 AND stat_type_id = $2;

-- name: DeletePerformanceStats :exec
DELETE FROM performance_stats
WHERE performance_id = $1;

-- name: GetTopPerformancesByStat :many
SELECT 
    ps.performance_id, ps.value,
    p.title AS performance_title, p.average_rating,
    pl.name AS player_name, pl.slug AS player_slug, pl.photo_url AS player_photo_url,
    m.title AS match_title, m.slug AS match_slug, m.utc_datetime AS match_utc_datetime
FROM performance_stats ps
JOIN performances p ON ps.performance_id = p.id
JOIN players pl ON p.player_id = pl.id
JOIN matches m ON p.match_id = m.id
WHERE ps.stat_type_id = $1
ORDER BY ps.value DESC, p.average_rating DESC
LIMIT $2 OFFSET $3;
