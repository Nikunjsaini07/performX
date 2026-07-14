-- name: GetMatchStats :many
SELECT 
    ms.id, ms.match_id, ms.team_id, ms.stat_type_id, ms.value,
    st.name AS stat_name, st.short_name AS stat_short_name, st.unit AS stat_unit, st.category AS stat_category
FROM match_stats ms
JOIN stat_types st ON ms.stat_type_id = st.id
WHERE ms.match_id = $1;

-- name: GetTeamMatchStats :many
SELECT 
    ms.id, ms.match_id, ms.team_id, ms.stat_type_id, ms.value,
    st.name AS stat_name, st.short_name AS stat_short_name, st.unit AS stat_unit, st.category AS stat_category
FROM match_stats ms
JOIN stat_types st ON ms.stat_type_id = st.id
WHERE ms.match_id = $1 AND ms.team_id = $2;

-- name: CreateMatchStat :one
INSERT INTO match_stats (
    match_id,
    team_id,
    stat_type_id,
    value
) VALUES (
    $1, $2, $3, $4
)
RETURNING id, match_id, team_id, stat_type_id, value, created_at;

-- name: CreateMatchStats :copyfrom
INSERT INTO match_stats (
    match_id,
    team_id,
    stat_type_id,
    value
) VALUES (
    $1, $2, $3, $4
);

-- name: UpdateMatchStat :one
UPDATE match_stats
SET value = $4
WHERE match_id = $1 AND team_id = $2 AND stat_type_id = $3
RETURNING id, match_id, team_id, stat_type_id, value, created_at;

-- name: DeleteMatchStat :exec
DELETE FROM match_stats
WHERE match_id = $1 AND team_id = $2 AND stat_type_id = $3;

-- name: DeleteMatchStats :exec
DELETE FROM match_stats
WHERE match_id = $1;

-- name: GetTopMatchesByStat :many
SELECT 
    ms.match_id, ms.team_id, ms.value,
    t.name AS team_name, t.logo_url AS team_logo_url,
    m.title AS match_title, m.slug AS match_slug, m.utc_datetime AS match_utc_datetime
FROM match_stats ms
JOIN teams t ON ms.team_id = t.id
JOIN matches m ON ms.match_id = m.id
WHERE ms.stat_type_id = $1
ORDER BY ms.value DESC
LIMIT $2 OFFSET $3;
