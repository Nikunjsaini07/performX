

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createMatch = `-- name: CreateMatch :one
INSERT INTO matches (
    home_team_id,
    away_team_id,
    title,
    slug,
    description,
    round,
    utc_datetime,
    venue,
    cover_image_url,
    home_score,
    away_score,
    home_penalty_score,
    away_penalty_score,
    tagline
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
)
RETURNING id, home_team_id, away_team_id, title, slug, description, round, utc_datetime, venue, cover_image_url, home_score, away_score, home_penalty_score, away_penalty_score, tagline, created_at, updated_at
`

type CreateMatchParams struct {
	HomeTeamID       pgtype.UUID        `db:"home_team_id" json:"home_team_id"`
	AwayTeamID       pgtype.UUID        `db:"away_team_id" json:"away_team_id"`
	Title            string             `db:"title" json:"title"`
	Slug             string             `db:"slug" json:"slug"`
	Description      pgtype.Text        `db:"description" json:"description"`
	Round            pgtype.Text        `db:"round" json:"round"`
	UtcDatetime      pgtype.Timestamptz `db:"utc_datetime" json:"utc_datetime"`
	Venue            pgtype.Text        `db:"venue" json:"venue"`
	CoverImageUrl    pgtype.Text        `db:"cover_image_url" json:"cover_image_url"`
	HomeScore        int32              `db:"home_score" json:"home_score"`
	AwayScore        int32              `db:"away_score" json:"away_score"`
	HomePenaltyScore pgtype.Int4        `db:"home_penalty_score" json:"home_penalty_score"`
	AwayPenaltyScore pgtype.Int4        `db:"away_penalty_score" json:"away_penalty_score"`
	Tagline          pgtype.Text        `db:"tagline" json:"tagline"`
}

type CreateMatchRow struct {
	ID               pgtype.UUID        `db:"id" json:"id"`
	HomeTeamID       pgtype.UUID        `db:"home_team_id" json:"home_team_id"`
	AwayTeamID       pgtype.UUID        `db:"away_team_id" json:"away_team_id"`
	Title            string             `db:"title" json:"title"`
	Slug             string             `db:"slug" json:"slug"`
	Description      pgtype.Text        `db:"description" json:"description"`
	Round            pgtype.Text        `db:"round" json:"round"`
	UtcDatetime      pgtype.Timestamptz `db:"utc_datetime" json:"utc_datetime"`
	Venue            pgtype.Text        `db:"venue" json:"venue"`
	CoverImageUrl    pgtype.Text        `db:"cover_image_url" json:"cover_image_url"`
	HomeScore        int32              `db:"home_score" json:"home_score"`
	AwayScore        int32              `db:"away_score" json:"away_score"`
	HomePenaltyScore pgtype.Int4        `db:"home_penalty_score" json:"home_penalty_score"`
	AwayPenaltyScore pgtype.Int4        `db:"away_penalty_score" json:"away_penalty_score"`
	Tagline          pgtype.Text        `db:"tagline" json:"tagline"`
	CreatedAt        pgtype.Timestamptz `db:"created_at" json:"created_at"`
	UpdatedAt        pgtype.Timestamptz `db:"updated_at" json:"updated_at"`
}

func (q *Queries) CreateMatch(ctx context.Context, arg CreateMatchParams) (CreateMatchRow, error) {
	row := q.db.QueryRow(ctx, createMatch,
		arg.HomeTeamID,
		arg.AwayTeamID,
		arg.Title,
		arg.Slug,
		arg.Description,
		arg.Round,
		arg.UtcDatetime,
		arg.Venue,
		arg.CoverImageUrl,
		arg.HomeScore,
		arg.AwayScore,
		arg.HomePenaltyScore,
		arg.AwayPenaltyScore,
		arg.Tagline,
	)
	var i CreateMatchRow
	err := row.Scan(
		&i.ID,
		&i.HomeTeamID,
		&i.AwayTeamID,
		&i.Title,
		&i.Slug,
		&i.Description,
		&i.Round,
		&i.UtcDatetime,
		&i.Venue,
		&i.CoverImageUrl,
		&i.HomeScore,
		&i.AwayScore,
		&i.HomePenaltyScore,
		&i.AwayPenaltyScore,
		&i.Tagline,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteMatch = `-- name: DeleteMatch :exec
DELETE FROM matches
WHERE id = $1
`

func (q *Queries) DeleteMatch(ctx context.Context, id pgtype.UUID) error {
	_, err := q.db.Exec(ctx, deleteMatch, id)
	return err
}

const getCompletedMatches = `-- name: GetCompletedMatches :many
SELECT 
    m.id, m.home_team_id, m.away_team_id, m.title, m.slug, m.tagline, m.utc_datetime, m.venue, m.home_score, m.away_score,
    ht.name AS home_team_name, ht.logo_url AS home_team_logo_url,
    at.name AS away_team_name, at.logo_url AS away_team_logo_url
FROM matches m
JOIN teams ht ON m.home_team_id = ht.id
JOIN teams at ON m.away_team_id = at.id
WHERE m.utc_datetime <= now()
ORDER BY m.utc_datetime DESC
LIMIT $1 OFFSET $2
`

type GetCompletedMatchesParams struct {
	Limit  int32 `db:"limit" json:"limit"`
	Offset int32 `db:"offset" json:"offset"`
}

type GetCompletedMatchesRow struct {
	ID              pgtype.UUID        `db:"id" json:"id"`
	HomeTeamID      pgtype.UUID        `db:"home_team_id" json:"home_team_id"`
	AwayTeamID      pgtype.UUID        `db:"away_team_id" json:"away_team_id"`
	Title           string             `db:"title" json:"title"`
	Slug            string             `db:"slug" json:"slug"`
	Tagline         pgtype.Text        `db:"tagline" json:"tagline"`
	UtcDatetime     pgtype.Timestamptz `db:"utc_datetime" json:"utc_datetime"`
	Venue           pgtype.Text        `db:"venue" json:"venue"`
	HomeScore       int32              `db:"home_score" json:"home_score"`
	AwayScore       int32              `db:"away_score" json:"away_score"`
	HomeTeamName    string             `db:"home_team_name" json:"home_team_name"`
	HomeTeamLogoUrl pgtype.Text        `db:"home_team_logo_url" json:"home_team_logo_url"`
	AwayTeamName    string             `db:"away_team_name" json:"away_team_name"`
	AwayTeamLogoUrl pgtype.Text        `db:"away_team_logo_url" json:"away_team_logo_url"`
}

func (q *Queries) GetCompletedMatches(ctx context.Context, arg GetCompletedMatchesParams) ([]GetCompletedMatchesRow, error) {
	rows, err := q.db.Query(ctx, getCompletedMatches, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetCompletedMatchesRow
	for rows.Next() {
		var i GetCompletedMatchesRow
		if err := rows.Scan(
			&i.ID,
			&i.HomeTeamID,
			&i.AwayTeamID,
			&i.Title,
			&i.Slug,
			&i.Tagline,
			&i.UtcDatetime,
			&i.Venue,
			&i.HomeScore,
			&i.AwayScore,
			&i.HomeTeamName,
			&i.HomeTeamLogoUrl,
			&i.AwayTeamName,
			&i.AwayTeamLogoUrl,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getMatchByID = `-- name: GetMatchByID :one
SELECT 
    m.id, m.home_team_id, m.away_team_id, m.title, m.slug, m.description, m.round, m.utc_datetime, m.venue, m.cover_image_url, m.home_score, m.away_score, m.home_penalty_score, m.away_penalty_score, m.tagline, m.created_at, m.updated_at,
    ht.name AS home_team_name, ht.short_name AS home_team_short_name, ht.logo_url AS home_team_logo_url,
    at.name AS away_team_name, at.short_name AS away_team_short_name, at.logo_url AS away_team_logo_url,
    m.average_rating,
    m.total_votes
FROM matches m
JOIN teams ht ON m.home_team_id = ht.id
JOIN teams at ON m.away_team_id = at.id
WHERE m.id = $1 LIMIT 1
`

type GetMatchByIDRow struct {
	ID                pgtype.UUID        `db:"id" json:"id"`
	HomeTeamID        pgtype.UUID        `db:"home_team_id" json:"home_team_id"`
	AwayTeamID        pgtype.UUID        `db:"away_team_id" json:"away_team_id"`
	Title             string             `db:"title" json:"title"`
	Slug              string             `db:"slug" json:"slug"`
	Description       pgtype.Text        `db:"description" json:"description"`
	Round             pgtype.Text        `db:"round" json:"round"`
	UtcDatetime       pgtype.Timestamptz `db:"utc_datetime" json:"utc_datetime"`
	Venue             pgtype.Text        `db:"venue" json:"venue"`
	CoverImageUrl     pgtype.Text        `db:"cover_image_url" json:"cover_image_url"`
	HomeScore         int32              `db:"home_score" json:"home_score"`
	AwayScore         int32              `db:"away_score" json:"away_score"`
	HomePenaltyScore  pgtype.Int4        `db:"home_penalty_score" json:"home_penalty_score"`
	AwayPenaltyScore  pgtype.Int4        `db:"away_penalty_score" json:"away_penalty_score"`
	Tagline           pgtype.Text        `db:"tagline" json:"tagline"`
	CreatedAt         pgtype.Timestamptz `db:"created_at" json:"created_at"`
	UpdatedAt         pgtype.Timestamptz `db:"updated_at" json:"updated_at"`
	HomeTeamName      string             `db:"home_team_name" json:"home_team_name"`
	HomeTeamShortName string             `db:"home_team_short_name" json:"home_team_short_name"`
	HomeTeamLogoUrl   pgtype.Text        `db:"home_team_logo_url" json:"home_team_logo_url"`
	AwayTeamName      string             `db:"away_team_name" json:"away_team_name"`
	AwayTeamShortName string             `db:"away_team_short_name" json:"away_team_short_name"`
	AwayTeamLogoUrl   pgtype.Text        `db:"away_team_logo_url" json:"away_team_logo_url"`
	AverageRating     pgtype.Numeric     `db:"average_rating" json:"average_rating"`
	TotalVotes        pgtype.Int4        `db:"total_votes" json:"total_votes"`
}

func (q *Queries) GetMatchByID(ctx context.Context, id pgtype.UUID) (GetMatchByIDRow, error) {
	row := q.db.QueryRow(ctx, getMatchByID, id)
	var i GetMatchByIDRow
	err := row.Scan(
		&i.ID,
		&i.HomeTeamID,
		&i.AwayTeamID,
		&i.Title,
		&i.Slug,
		&i.Description,
		&i.Round,
		&i.UtcDatetime,
		&i.Venue,
		&i.CoverImageUrl,
		&i.HomeScore,
		&i.AwayScore,
		&i.HomePenaltyScore,
		&i.AwayPenaltyScore,
		&i.Tagline,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.HomeTeamName,
		&i.HomeTeamShortName,
		&i.HomeTeamLogoUrl,
		&i.AwayTeamName,
		&i.AwayTeamShortName,
		&i.AwayTeamLogoUrl,
		&i.AverageRating,
		&i.TotalVotes,
	)
	return i, err
}

const getMatchBySlug = `-- name: GetMatchBySlug :one
SELECT 
    m.id, m.home_team_id, m.away_team_id, m.title, m.slug, m.description, m.round, m.utc_datetime, m.venue, m.cover_image_url, m.home_score, m.away_score, m.home_penalty_score, m.away_penalty_score, m.tagline, m.created_at, m.updated_at,
    ht.name AS home_team_name, ht.short_name AS home_team_short_name, ht.logo_url AS home_team_logo_url,
    at.name AS away_team_name, at.short_name AS away_team_short_name, at.logo_url AS away_team_logo_url,
    m.average_rating,
    m.total_votes
FROM matches m
JOIN teams ht ON m.home_team_id = ht.id
JOIN teams at ON m.away_team_id = at.id
WHERE m.slug = $1 LIMIT 1
`

type GetMatchBySlugRow struct {
	ID                pgtype.UUID        `db:"id" json:"id"`
	HomeTeamID        pgtype.UUID        `db:"home_team_id" json:"home_team_id"`
	AwayTeamID        pgtype.UUID        `db:"away_team_id" json:"away_team_id"`
	Title             string             `db:"title" json:"title"`
	Slug              string             `db:"slug" json:"slug"`
	Description       pgtype.Text        `db:"description" json:"description"`
	Round             pgtype.Text        `db:"round" json:"round"`
	UtcDatetime       pgtype.Timestamptz `db:"utc_datetime" json:"utc_datetime"`
	Venue             pgtype.Text        `db:"venue" json:"venue"`
	CoverImageUrl     pgtype.Text        `db:"cover_image_url" json:"cover_image_url"`
	HomeScore         int32              `db:"home_score" json:"home_score"`
	AwayScore         int32              `db:"away_score" json:"away_score"`
	HomePenaltyScore  pgtype.Int4        `db:"home_penalty_score" json:"home_penalty_score"`
	AwayPenaltyScore  pgtype.Int4        `db:"away_penalty_score" json:"away_penalty_score"`
	Tagline           pgtype.Text        `db:"tagline" json:"tagline"`
	CreatedAt         pgtype.Timestamptz `db:"created_at" json:"created_at"`
	UpdatedAt         pgtype.Timestamptz `db:"updated_at" json:"updated_at"`
	HomeTeamName      string             `db:"home_team_name" json:"home_team_name"`
	HomeTeamShortName string             `db:"home_team_short_name" json:"home_team_short_name"`
	HomeTeamLogoUrl   pgtype.Text        `db:"home_team_logo_url" json:"home_team_logo_url"`
	AwayTeamName      string             `db:"away_team_name" json:"away_team_name"`
	AwayTeamShortName string             `db:"away_team_short_name" json:"away_team_short_name"`
	AwayTeamLogoUrl   pgtype.Text        `db:"away_team_logo_url" json:"away_team_logo_url"`
	AverageRating     pgtype.Numeric     `db:"average_rating" json:"average_rating"`
	TotalVotes        pgtype.Int4        `db:"total_votes" json:"total_votes"`
}

func (q *Queries) GetMatchBySlug(ctx context.Context, slug string) (GetMatchBySlugRow, error) {
	row := q.db.QueryRow(ctx, getMatchBySlug, slug)
	var i GetMatchBySlugRow
	err := row.Scan(
		&i.ID,
		&i.HomeTeamID,
		&i.AwayTeamID,
		&i.Title,
		&i.Slug,
		&i.Description,
		&i.Round,
		&i.UtcDatetime,
		&i.Venue,
		&i.CoverImageUrl,
		&i.HomeScore,
		&i.AwayScore,
		&i.HomePenaltyScore,
		&i.AwayPenaltyScore,
		&i.Tagline,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.HomeTeamName,
		&i.HomeTeamShortName,
		&i.HomeTeamLogoUrl,
		&i.AwayTeamName,
		&i.AwayTeamShortName,
		&i.AwayTeamLogoUrl,
		&i.AverageRating,
		&i.TotalVotes,
	)
	return i, err
}

const getMatchPerformances = `-- name: GetMatchPerformances :many
SELECT 
    p.id, p.match_id, p.player_id, p.player_team_id, p.title, p.cover_image_url, p.jersey_number, p.is_starter, p.captain, p.minutes_played, p.average_rating, p.slug, p.tagline,
    pl.name AS player_name, pl.slug AS player_slug, pl.photo_url AS player_photo_url,
    t.name AS team_name, t.logo_url AS team_logo_url, t.flag_emoji,
    COALESCE((SELECT sv.value FROM performance_stats sv JOIN stat_types stt ON stt.id = sv.stat_type_id WHERE sv.performance_id = p.id AND stt.name = 'goals' LIMIT 1), 0)::int AS goals,
    COALESCE((SELECT sv.value FROM performance_stats sv JOIN stat_types stt ON stt.id = sv.stat_type_id WHERE sv.performance_id = p.id AND stt.name = 'assists' LIMIT 1), 0)::int AS assists,
    0::int AS dribbles,
    0::numeric(5,2) AS passes_accuracy,
    p.average_rating::numeric(3,1) AS average_rating,
    p.total_votes::bigint AS total_votes
FROM performances p
JOIN players pl ON p.player_id = pl.id
JOIN player_teams pt ON p.player_team_id = pt.id
JOIN teams t ON pt.team_id = t.id
WHERE p.match_id = $1
ORDER BY p.is_starter DESC, p.minutes_played DESC
`

type GetMatchPerformancesRow struct {
	ID              pgtype.UUID    `db:"id" json:"id"`
	MatchID         pgtype.UUID    `db:"match_id" json:"match_id"`
	PlayerID        pgtype.UUID    `db:"player_id" json:"player_id"`
	PlayerTeamID    pgtype.UUID    `db:"player_team_id" json:"player_team_id"`
	Title           string         `db:"title" json:"title"`
	CoverImageUrl   pgtype.Text    `db:"cover_image_url" json:"cover_image_url"`
	JerseyNumber    pgtype.Int4    `db:"jersey_number" json:"jersey_number"`
	IsStarter       bool           `db:"is_starter" json:"is_starter"`
	Captain         bool           `db:"captain" json:"captain"`
	MinutesPlayed   int32          `db:"minutes_played" json:"minutes_played"`
	AverageRating   pgtype.Numeric `db:"average_rating" json:"average_rating"`
	Slug            pgtype.Text    `db:"slug" json:"slug"`
	Tagline         pgtype.Text    `db:"tagline" json:"tagline"`
	PlayerName      string         `db:"player_name" json:"player_name"`
	PlayerSlug      string         `db:"player_slug" json:"player_slug"`
	PlayerPhotoUrl  pgtype.Text    `db:"player_photo_url" json:"player_photo_url"`
	TeamName        string         `db:"team_name" json:"team_name"`
	TeamLogoUrl     pgtype.Text    `db:"team_logo_url" json:"team_logo_url"`
	FlagEmoji       pgtype.Text    `db:"flag_emoji" json:"flag_emoji"`
	Goals           int32          `db:"goals" json:"goals"`
	Assists         int32          `db:"assists" json:"assists"`
	Dribbles        int32          `db:"dribbles" json:"dribbles"`
	PassesAccuracy  pgtype.Numeric `db:"passes_accuracy" json:"passes_accuracy"`
	AverageRating_2 pgtype.Numeric `db:"average_rating_2" json:"average_rating_2"`
	TotalVotes      int64          `db:"total_votes" json:"total_votes"`
}

func (q *Queries) GetMatchPerformances(ctx context.Context, matchID pgtype.UUID) ([]GetMatchPerformancesRow, error) {
	rows, err := q.db.Query(ctx, getMatchPerformances, matchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetMatchPerformancesRow
	for rows.Next() {
		var i GetMatchPerformancesRow
		if err := rows.Scan(
			&i.ID,
			&i.MatchID,
			&i.PlayerID,
			&i.PlayerTeamID,
			&i.Title,
			&i.CoverImageUrl,
			&i.JerseyNumber,
			&i.IsStarter,
			&i.Captain,
			&i.MinutesPlayed,
			&i.AverageRating,
			&i.Slug,
			&i.Tagline,
			&i.PlayerName,
			&i.PlayerSlug,
			&i.PlayerPhotoUrl,
			&i.TeamName,
			&i.TeamLogoUrl,
			&i.FlagEmoji,
			&i.Goals,
			&i.Assists,
			&i.Dribbles,
			&i.PassesAccuracy,
			&i.AverageRating_2,
			&i.TotalVotes,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getMatchesByDate = `-- name: GetMatchesByDate :many
SELECT 
    m.id, m.home_team_id, m.away_team_id, m.title, m.slug, m.tagline, m.utc_datetime, m.home_score, m.away_score,
    ht.name AS home_team_name, ht.logo_url AS home_team_logo_url,
    at.name AS away_team_name, at.logo_url AS away_team_logo_url
FROM matches m
JOIN teams ht ON m.home_team_id = ht.id
JOIN teams at ON m.away_team_id = at.id
WHERE m.utc_datetime::date = $1::date
ORDER BY m.utc_datetime DESC
`

type GetMatchesByDateRow struct {
	ID              pgtype.UUID        `db:"id" json:"id"`
	HomeTeamID      pgtype.UUID        `db:"home_team_id" json:"home_team_id"`
	AwayTeamID      pgtype.UUID        `db:"away_team_id" json:"away_team_id"`
	Title           string             `db:"title" json:"title"`
	Slug            string             `db:"slug" json:"slug"`
	Tagline         pgtype.Text        `db:"tagline" json:"tagline"`
	UtcDatetime     pgtype.Timestamptz `db:"utc_datetime" json:"utc_datetime"`
	HomeScore       int32              `db:"home_score" json:"home_score"`
	AwayScore       int32              `db:"away_score" json:"away_score"`
	HomeTeamName    string             `db:"home_team_name" json:"home_team_name"`
	HomeTeamLogoUrl pgtype.Text        `db:"home_team_logo_url" json:"home_team_logo_url"`
	AwayTeamName    string             `db:"away_team_name" json:"away_team_name"`
	AwayTeamLogoUrl pgtype.Text        `db:"away_team_logo_url" json:"away_team_logo_url"`
}

func (q *Queries) GetMatchesByDate(ctx context.Context, dollar_1 pgtype.Date) ([]GetMatchesByDateRow, error) {
	rows, err := q.db.Query(ctx, getMatchesByDate, dollar_1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetMatchesByDateRow
	for rows.Next() {
		var i GetMatchesByDateRow
		if err := rows.Scan(
			&i.ID,
			&i.HomeTeamID,
			&i.AwayTeamID,
			&i.Title,
			&i.Slug,
			&i.Tagline,
			&i.UtcDatetime,
			&i.HomeScore,
			&i.AwayScore,
			&i.HomeTeamName,
			&i.HomeTeamLogoUrl,
			&i.AwayTeamName,
			&i.AwayTeamLogoUrl,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getMatchesByTeam = `-- name: GetMatchesByTeam :many
SELECT 
    m.id, m.home_team_id, m.away_team_id, m.title, m.slug, m.tagline, m.utc_datetime, m.home_score, m.away_score,
    ht.name AS home_team_name, ht.logo_url AS home_team_logo_url,
    at.name AS away_team_name, at.logo_url AS away_team_logo_url
FROM matches m
JOIN teams ht ON m.home_team_id = ht.id
JOIN teams at ON m.away_team_id = at.id
WHERE m.home_team_id = $1 OR m.away_team_id = $1
ORDER BY m.utc_datetime DESC
LIMIT $2 OFFSET $3
`

type GetMatchesByTeamParams struct {
	HomeTeamID pgtype.UUID `db:"home_team_id" json:"home_team_id"`
	Limit      int32       `db:"limit" json:"limit"`
	Offset     int32       `db:"offset" json:"offset"`
}

type GetMatchesByTeamRow struct {
	ID              pgtype.UUID        `db:"id" json:"id"`
	HomeTeamID      pgtype.UUID        `db:"home_team_id" json:"home_team_id"`
	AwayTeamID      pgtype.UUID        `db:"away_team_id" json:"away_team_id"`
	Title           string             `db:"title" json:"title"`
	Slug            string             `db:"slug" json:"slug"`
	Tagline         pgtype.Text        `db:"tagline" json:"tagline"`
	UtcDatetime     pgtype.Timestamptz `db:"utc_datetime" json:"utc_datetime"`
	HomeScore       int32              `db:"home_score" json:"home_score"`
	AwayScore       int32              `db:"away_score" json:"away_score"`
	HomeTeamName    string             `db:"home_team_name" json:"home_team_name"`
	HomeTeamLogoUrl pgtype.Text        `db:"home_team_logo_url" json:"home_team_logo_url"`
	AwayTeamName    string             `db:"away_team_name" json:"away_team_name"`
	AwayTeamLogoUrl pgtype.Text        `db:"away_team_logo_url" json:"away_team_logo_url"`
}

func (q *Queries) GetMatchesByTeam(ctx context.Context, arg GetMatchesByTeamParams) ([]GetMatchesByTeamRow, error) {
	rows, err := q.db.Query(ctx, getMatchesByTeam, arg.HomeTeamID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetMatchesByTeamRow
	for rows.Next() {
		var i GetMatchesByTeamRow
		if err := rows.Scan(
			&i.ID,
			&i.HomeTeamID,
			&i.AwayTeamID,
			&i.Title,
			&i.Slug,
			&i.Tagline,
			&i.UtcDatetime,
			&i.HomeScore,
			&i.AwayScore,
			&i.HomeTeamName,
			&i.HomeTeamLogoUrl,
			&i.AwayTeamName,
			&i.AwayTeamLogoUrl,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getOverviewStats = `-- name: GetOverviewStats :one
SELECT
    (SELECT COUNT(*) FROM matches)::bigint AS match_count,
    (SELECT COUNT(*) FROM performances)::bigint AS performance_count,
    (SELECT COUNT(*) FROM teams)::bigint AS team_count,
    ((SELECT COUNT(*) FROM match_ratings) + (SELECT COUNT(*) FROM performance_ratings))::bigint AS rating_count,
    ((SELECT COUNT(*) FROM match_reviews) + (SELECT COUNT(*) FROM performance_reviews))::bigint AS review_count
`

type GetOverviewStatsRow struct {
	MatchCount       int64 `db:"match_count" json:"match_count"`
	PerformanceCount int64 `db:"performance_count" json:"performance_count"`
	TeamCount        int64 `db:"team_count" json:"team_count"`
	RatingCount      int64 `db:"rating_count" json:"rating_count"`
	ReviewCount      int64 `db:"review_count" json:"review_count"`
}

func (q *Queries) GetOverviewStats(ctx context.Context) (GetOverviewStatsRow, error) {
	row := q.db.QueryRow(ctx, getOverviewStats)
	var i GetOverviewStatsRow
	err := row.Scan(
		&i.MatchCount,
		&i.PerformanceCount,
		&i.TeamCount,
		&i.RatingCount,
		&i.ReviewCount,
	)
	return i, err
}

const getRelatedMatches = `-- name: GetRelatedMatches :many
SELECT 
    m.id, m.home_team_id, m.away_team_id, m.title, m.slug, m.tagline, m.utc_datetime, m.home_score, m.away_score,
    ht.name AS home_team_name, ht.logo_url AS home_team_logo_url,
    at.name AS away_team_name, at.logo_url AS away_team_logo_url
FROM matches m
JOIN teams ht ON m.home_team_id = ht.id
JOIN teams at ON m.away_team_id = at.id
WHERE m.id <> $1
ORDER BY m.utc_datetime DESC
LIMIT $2
`

type GetRelatedMatchesParams struct {
	ID    pgtype.UUID `db:"id" json:"id"`
	Limit int32       `db:"limit" json:"limit"`
}

type GetRelatedMatchesRow struct {
	ID              pgtype.UUID        `db:"id" json:"id"`
	HomeTeamID      pgtype.UUID        `db:"home_team_id" json:"home_team_id"`
	AwayTeamID      pgtype.UUID        `db:"away_team_id" json:"away_team_id"`
	Title           string             `db:"title" json:"title"`
	Slug            string             `db:"slug" json:"slug"`
	Tagline         pgtype.Text        `db:"tagline" json:"tagline"`
	UtcDatetime     pgtype.Timestamptz `db:"utc_datetime" json:"utc_datetime"`
	HomeScore       int32              `db:"home_score" json:"home_score"`
	AwayScore       int32              `db:"away_score" json:"away_score"`
	HomeTeamName    string             `db:"home_team_name" json:"home_team_name"`
	HomeTeamLogoUrl pgtype.Text        `db:"home_team_logo_url" json:"home_team_logo_url"`
	AwayTeamName    string             `db:"away_team_name" json:"away_team_name"`
	AwayTeamLogoUrl pgtype.Text        `db:"away_team_logo_url" json:"away_team_logo_url"`
}

func (q *Queries) GetRelatedMatches(ctx context.Context, arg GetRelatedMatchesParams) ([]GetRelatedMatchesRow, error) {
	rows, err := q.db.Query(ctx, getRelatedMatches, arg.ID, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetRelatedMatchesRow
	for rows.Next() {
		var i GetRelatedMatchesRow
		if err := rows.Scan(
			&i.ID,
			&i.HomeTeamID,
			&i.AwayTeamID,
			&i.Title,
			&i.Slug,
			&i.Tagline,
			&i.UtcDatetime,
			&i.HomeScore,
			&i.AwayScore,
			&i.HomeTeamName,
			&i.HomeTeamLogoUrl,
			&i.AwayTeamName,
			&i.AwayTeamLogoUrl,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUpcomingMatches = `-- name: GetUpcomingMatches :many
SELECT 
    m.id, m.home_team_id, m.away_team_id, m.title, m.slug, m.tagline, m.utc_datetime, m.venue,
    ht.name AS home_team_name, ht.logo_url AS home_team_logo_url,
    at.name AS away_team_name, at.logo_url AS away_team_logo_url
FROM matches m
JOIN teams ht ON m.home_team_id = ht.id
JOIN teams at ON m.away_team_id = at.id
WHERE m.utc_datetime > now()
ORDER BY m.utc_datetime ASC
LIMIT $1 OFFSET $2
`

type GetUpcomingMatchesParams struct {
	Limit  int32 `db:"limit" json:"limit"`
	Offset int32 `db:"offset" json:"offset"`
}

type GetUpcomingMatchesRow struct {
	ID              pgtype.UUID        `db:"id" json:"id"`
	HomeTeamID      pgtype.UUID        `db:"home_team_id" json:"home_team_id"`
	AwayTeamID      pgtype.UUID        `db:"away_team_id" json:"away_team_id"`
	Title           string             `db:"title" json:"title"`
	Slug            string             `db:"slug" json:"slug"`
	Tagline         pgtype.Text        `db:"tagline" json:"tagline"`
	UtcDatetime     pgtype.Timestamptz `db:"utc_datetime" json:"utc_datetime"`
	Venue           pgtype.Text        `db:"venue" json:"venue"`
	HomeTeamName    string             `db:"home_team_name" json:"home_team_name"`
	HomeTeamLogoUrl pgtype.Text        `db:"home_team_logo_url" json:"home_team_logo_url"`
	AwayTeamName    string             `db:"away_team_name" json:"away_team_name"`
	AwayTeamLogoUrl pgtype.Text        `db:"away_team_logo_url" json:"away_team_logo_url"`
}

func (q *Queries) GetUpcomingMatches(ctx context.Context, arg GetUpcomingMatchesParams) ([]GetUpcomingMatchesRow, error) {
	rows, err := q.db.Query(ctx, getUpcomingMatches, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetUpcomingMatchesRow
	for rows.Next() {
		var i GetUpcomingMatchesRow
		if err := rows.Scan(
			&i.ID,
			&i.HomeTeamID,
			&i.AwayTeamID,
			&i.Title,
			&i.Slug,
			&i.Tagline,
			&i.UtcDatetime,
			&i.Venue,
			&i.HomeTeamName,
			&i.HomeTeamLogoUrl,
			&i.AwayTeamName,
			&i.AwayTeamLogoUrl,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listMatches = `-- name: ListMatches :many
SELECT 
    m.id, m.home_team_id, m.away_team_id, m.title, m.slug, m.description, m.round, m.utc_datetime, m.venue, m.cover_image_url, m.home_score, m.away_score, m.home_penalty_score, m.away_penalty_score, m.tagline, m.created_at, m.updated_at,
    ht.name AS home_team_name, ht.short_name AS home_team_short_name, ht.logo_url AS home_team_logo_url,
    at.name AS away_team_name, at.short_name AS away_team_short_name, at.logo_url AS away_team_logo_url,
    m.average_rating,
    m.total_votes
FROM matches m
JOIN teams ht ON m.home_team_id = ht.id
JOIN teams at ON m.away_team_id = at.id
ORDER BY m.utc_datetime DESC
LIMIT $1 OFFSET $2
`

type ListMatchesParams struct {
	Limit  int32 `db:"limit" json:"limit"`
	Offset int32 `db:"offset" json:"offset"`
}

type ListMatchesRow struct {
	ID                pgtype.UUID        `db:"id" json:"id"`
	HomeTeamID        pgtype.UUID        `db:"home_team_id" json:"home_team_id"`
	AwayTeamID        pgtype.UUID        `db:"away_team_id" json:"away_team_id"`
	Title             string             `db:"title" json:"title"`
	Slug              string             `db:"slug" json:"slug"`
	Description       pgtype.Text        `db:"description" json:"description"`
	Round             pgtype.Text        `db:"round" json:"round"`
	UtcDatetime       pgtype.Timestamptz `db:"utc_datetime" json:"utc_datetime"`
	Venue             pgtype.Text        `db:"venue" json:"venue"`
	CoverImageUrl     pgtype.Text        `db:"cover_image_url" json:"cover_image_url"`
	HomeScore         int32              `db:"home_score" json:"home_score"`
	AwayScore         int32              `db:"away_score" json:"away_score"`
	HomePenaltyScore  pgtype.Int4        `db:"home_penalty_score" json:"home_penalty_score"`
	AwayPenaltyScore  pgtype.Int4        `db:"away_penalty_score" json:"away_penalty_score"`
	Tagline           pgtype.Text        `db:"tagline" json:"tagline"`
	CreatedAt         pgtype.Timestamptz `db:"created_at" json:"created_at"`
	UpdatedAt         pgtype.Timestamptz `db:"updated_at" json:"updated_at"`
	HomeTeamName      string             `db:"home_team_name" json:"home_team_name"`
	HomeTeamShortName string             `db:"home_team_short_name" json:"home_team_short_name"`
	HomeTeamLogoUrl   pgtype.Text        `db:"home_team_logo_url" json:"home_team_logo_url"`
	AwayTeamName      string             `db:"away_team_name" json:"away_team_name"`
	AwayTeamShortName string             `db:"away_team_short_name" json:"away_team_short_name"`
	AwayTeamLogoUrl   pgtype.Text        `db:"away_team_logo_url" json:"away_team_logo_url"`
	AverageRating     pgtype.Numeric     `db:"average_rating" json:"average_rating"`
	TotalVotes        pgtype.Int4        `db:"total_votes" json:"total_votes"`
}

func (q *Queries) ListMatches(ctx context.Context, arg ListMatchesParams) ([]ListMatchesRow, error) {
	rows, err := q.db.Query(ctx, listMatches, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListMatchesRow
	for rows.Next() {
		var i ListMatchesRow
		if err := rows.Scan(
			&i.ID,
			&i.HomeTeamID,
			&i.AwayTeamID,
			&i.Title,
			&i.Slug,
			&i.Description,
			&i.Round,
			&i.UtcDatetime,
			&i.Venue,
			&i.CoverImageUrl,
			&i.HomeScore,
			&i.AwayScore,
			&i.HomePenaltyScore,
			&i.AwayPenaltyScore,
			&i.Tagline,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.HomeTeamName,
			&i.HomeTeamShortName,
			&i.HomeTeamLogoUrl,
			&i.AwayTeamName,
			&i.AwayTeamShortName,
			&i.AwayTeamLogoUrl,
			&i.AverageRating,
			&i.TotalVotes,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const searchMatches = `-- name: SearchMatches :many
SELECT 
    m.id, m.home_team_id, m.away_team_id, m.title, m.slug, m.tagline, m.utc_datetime, m.venue, m.home_score, m.away_score,
    ht.name AS home_team_name, ht.logo_url AS home_team_logo_url,
    at.name AS away_team_name, at.logo_url AS away_team_logo_url
FROM matches m
JOIN teams ht ON m.home_team_id = ht.id
JOIN teams at ON m.away_team_id = at.id
WHERE m.title ILIKE $1 OR ht.name ILIKE $1 OR at.name ILIKE $1
ORDER BY m.utc_datetime DESC
LIMIT $2 OFFSET $3
`

type SearchMatchesParams struct {
	Title  string `db:"title" json:"title"`
	Limit  int32  `db:"limit" json:"limit"`
	Offset int32  `db:"offset" json:"offset"`
}

type SearchMatchesRow struct {
	ID              pgtype.UUID        `db:"id" json:"id"`
	HomeTeamID      pgtype.UUID        `db:"home_team_id" json:"home_team_id"`
	AwayTeamID      pgtype.UUID        `db:"away_team_id" json:"away_team_id"`
	Title           string             `db:"title" json:"title"`
	Slug            string             `db:"slug" json:"slug"`
	Tagline         pgtype.Text        `db:"tagline" json:"tagline"`
	UtcDatetime     pgtype.Timestamptz `db:"utc_datetime" json:"utc_datetime"`
	Venue           pgtype.Text        `db:"venue" json:"venue"`
	HomeScore       int32              `db:"home_score" json:"home_score"`
	AwayScore       int32              `db:"away_score" json:"away_score"`
	HomeTeamName    string             `db:"home_team_name" json:"home_team_name"`
	HomeTeamLogoUrl pgtype.Text        `db:"home_team_logo_url" json:"home_team_logo_url"`
	AwayTeamName    string             `db:"away_team_name" json:"away_team_name"`
	AwayTeamLogoUrl pgtype.Text        `db:"away_team_logo_url" json:"away_team_logo_url"`
}

func (q *Queries) SearchMatches(ctx context.Context, arg SearchMatchesParams) ([]SearchMatchesRow, error) {
	rows, err := q.db.Query(ctx, searchMatches, arg.Title, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SearchMatchesRow
	for rows.Next() {
		var i SearchMatchesRow
		if err := rows.Scan(
			&i.ID,
			&i.HomeTeamID,
			&i.AwayTeamID,
			&i.Title,
			&i.Slug,
			&i.Tagline,
			&i.UtcDatetime,
			&i.Venue,
			&i.HomeScore,
			&i.AwayScore,
			&i.HomeTeamName,
			&i.HomeTeamLogoUrl,
			&i.AwayTeamName,
			&i.AwayTeamLogoUrl,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateMatch = `-- name: UpdateMatch :one
UPDATE matches
SET
    home_team_id = COALESCE($2, home_team_id),
    away_team_id = COALESCE($3, away_team_id),
    title = COALESCE($4, title),
    slug = COALESCE($5, slug),
    description = COALESCE($6, description),
    round = COALESCE($7, round),
    utc_datetime = COALESCE($8, utc_datetime),
    venue = COALESCE($9, venue),
    cover_image_url = COALESCE($10, cover_image_url),
    home_score = COALESCE($11, home_score),
    away_score = COALESCE($12, away_score),
    home_penalty_score = COALESCE($13, home_penalty_score),
    away_penalty_score = COALESCE($14, away_penalty_score),
    tagline = COALESCE($15, tagline),
    updated_at = now()
WHERE id = $1
RETURNING id, home_team_id, away_team_id, title, slug, description, round, utc_datetime, venue, cover_image_url, home_score, away_score, home_penalty_score, away_penalty_score, tagline, created_at, updated_at
`

type UpdateMatchParams struct {
	ID               pgtype.UUID        `db:"id" json:"id"`
	HomeTeamID       pgtype.UUID        `db:"home_team_id" json:"home_team_id"`
	AwayTeamID       pgtype.UUID        `db:"away_team_id" json:"away_team_id"`
	Title            string             `db:"title" json:"title"`
	Slug             string             `db:"slug" json:"slug"`
	Description      pgtype.Text        `db:"description" json:"description"`
	Round            pgtype.Text        `db:"round" json:"round"`
	UtcDatetime      pgtype.Timestamptz `db:"utc_datetime" json:"utc_datetime"`
	Venue            pgtype.Text        `db:"venue" json:"venue"`
	CoverImageUrl    pgtype.Text        `db:"cover_image_url" json:"cover_image_url"`
	HomeScore        int32              `db:"home_score" json:"home_score"`
	AwayScore        int32              `db:"away_score" json:"away_score"`
	HomePenaltyScore pgtype.Int4        `db:"home_penalty_score" json:"home_penalty_score"`
	AwayPenaltyScore pgtype.Int4        `db:"away_penalty_score" json:"away_penalty_score"`
	Tagline          pgtype.Text        `db:"tagline" json:"tagline"`
}

type UpdateMatchRow struct {
	ID               pgtype.UUID        `db:"id" json:"id"`
	HomeTeamID       pgtype.UUID        `db:"home_team_id" json:"home_team_id"`
	AwayTeamID       pgtype.UUID        `db:"away_team_id" json:"away_team_id"`
	Title            string             `db:"title" json:"title"`
	Slug             string             `db:"slug" json:"slug"`
	Description      pgtype.Text        `db:"description" json:"description"`
	Round            pgtype.Text        `db:"round" json:"round"`
	UtcDatetime      pgtype.Timestamptz `db:"utc_datetime" json:"utc_datetime"`
	Venue            pgtype.Text        `db:"venue" json:"venue"`
	CoverImageUrl    pgtype.Text        `db:"cover_image_url" json:"cover_image_url"`
	HomeScore        int32              `db:"home_score" json:"home_score"`
	AwayScore        int32              `db:"away_score" json:"away_score"`
	HomePenaltyScore pgtype.Int4        `db:"home_penalty_score" json:"home_penalty_score"`
	AwayPenaltyScore pgtype.Int4        `db:"away_penalty_score" json:"away_penalty_score"`
	Tagline          pgtype.Text        `db:"tagline" json:"tagline"`
	CreatedAt        pgtype.Timestamptz `db:"created_at" json:"created_at"`
	UpdatedAt        pgtype.Timestamptz `db:"updated_at" json:"updated_at"`
}

func (q *Queries) UpdateMatch(ctx context.Context, arg UpdateMatchParams) (UpdateMatchRow, error) {
	row := q.db.QueryRow(ctx, updateMatch,
		arg.ID,
		arg.HomeTeamID,
		arg.AwayTeamID,
		arg.Title,
		arg.Slug,
		arg.Description,
		arg.Round,
		arg.UtcDatetime,
		arg.Venue,
		arg.CoverImageUrl,
		arg.HomeScore,
		arg.AwayScore,
		arg.HomePenaltyScore,
		arg.AwayPenaltyScore,
		arg.Tagline,
	)
	var i UpdateMatchRow
	err := row.Scan(
		&i.ID,
		&i.HomeTeamID,
		&i.AwayTeamID,
		&i.Title,
		&i.Slug,
		&i.Description,
		&i.Round,
		&i.UtcDatetime,
		&i.Venue,
		&i.CoverImageUrl,
		&i.HomeScore,
		&i.AwayScore,
		&i.HomePenaltyScore,
		&i.AwayPenaltyScore,
		&i.Tagline,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
