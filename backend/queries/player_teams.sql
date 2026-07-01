-- name: JoinPlayerToTeam :one
INSERT INTO player_teams (
    player_id,
    team_id,
    jersey_number,
    start_date,
    end_date,
    is_active
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING id, player_id, team_id, jersey_number, start_date, end_date, is_active;

-- name: LeavePlayerTeam :one
UPDATE player_teams
SET 
    end_date = $2,
    is_active = FALSE
WHERE id = $1
RETURNING id, player_id, team_id, jersey_number, start_date, end_date, is_active;

-- name: GetPlayerCareer :many
SELECT 
    pt.id AS player_team_id, pt.jersey_number, pt.start_date, pt.end_date, pt.is_active,
    t.id AS team_id, t.name AS team_name, t.short_name AS team_short_name, t.logo_url AS team_logo_url, t.type AS team_type
FROM player_teams pt
JOIN teams t ON pt.team_id = t.id
WHERE pt.player_id = $1
ORDER BY pt.start_date DESC;

-- name: GetCurrentTeam :many
SELECT 
    pt.id AS player_team_id, pt.jersey_number, pt.start_date, pt.is_active,
    t.id AS team_id, t.name AS team_name, t.short_name AS team_short_name, t.logo_url AS team_logo_url, t.type AS team_type
FROM player_teams pt
JOIN teams t ON pt.team_id = t.id
WHERE pt.player_id = $1 AND pt.is_active = TRUE;

-- name: GetTeamPlayers :many
SELECT DISTINCT
    p.id, p.name, p.slug, p.known_as, p.photo_url, p.shirt_name,
    pt.jersey_number, pt.start_date, pt.end_date, pt.is_active
FROM player_teams pt
JOIN players p ON pt.player_id = p.id
WHERE pt.team_id = $1
ORDER BY pt.is_active DESC, p.name;

-- name: GetPlayerTeams :many
SELECT 
    pt.id, pt.player_id, pt.team_id, pt.jersey_number, pt.start_date, pt.end_date, pt.is_active,
    t.name AS team_name, t.logo_url AS team_logo_url
FROM player_teams pt
JOIN teams t ON pt.team_id = t.id
WHERE pt.player_id = $1
ORDER BY pt.start_date DESC;

-- name: UpdateJerseyNumber :one
UPDATE player_teams
SET jersey_number = $2
WHERE id = $1
RETURNING id, player_id, team_id, jersey_number, start_date, end_date, is_active;

-- name: UpdateSquadNumber :one
UPDATE player_teams
SET jersey_number = $2
WHERE id = $1
RETURNING id, player_id, team_id, jersey_number, start_date, end_date, is_active;

-- name: UpdatePlayerTeam :one
UPDATE player_teams
SET 
    jersey_number = COALESCE($2, jersey_number),
    start_date = COALESCE($3, start_date),
    end_date = COALESCE($4, end_date),
    is_active = COALESCE($5, is_active)
WHERE id = $1
RETURNING id, player_id, team_id, jersey_number, start_date, end_date, is_active;
