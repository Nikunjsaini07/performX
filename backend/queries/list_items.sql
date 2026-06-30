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
