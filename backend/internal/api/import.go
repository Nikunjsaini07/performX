package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/Nikunjsaini07/performx/backend/internal/db"
)

type ImportHandler struct {
	Queries *db.Queries
}

type ImportMatchRequest struct {
	HomeTeamName    string              `json:"home_team_name"`
	HomeTeamLogoURL string              `json:"home_team_logo_url"`
	AwayTeamName    string              `json:"away_team_name"`
	AwayTeamLogoURL string              `json:"away_team_logo_url"`
	Title           string              `json:"title"`
	UtcDatetime     string              `json:"utc_datetime"`
	Venue           string              `json:"venue"`
	HomeScore       int32               `json:"home_score"`
	AwayScore       int32               `json:"away_score"`
	CoverImageURL   string              `json:"cover_image_url"`
	Performances    []ImportPerformance `json:"performances"`
}

type ImportPerformance struct {
	PlayerName   string  `json:"player_name"`
	TeamName     string  `json:"team_name"`
	Rating       float64 `json:"rating"`
	JerseyNumber int32   `json:"jersey_number"`
	PhotoURL     string  `json:"photo_url"`
	Goals        int32   `json:"goals"`
	Assists      int32   `json:"assists"`
}

func makeSlug(s string) string {
	s = strings.ToLower(s)
	reg, _ := regexp.Compile("[^a-z0-9]+")
	s = reg.ReplaceAllString(s, "-")
	return strings.Trim(s, "-")
}

// POST /admin-gate-performx/import-match
func (h *ImportHandler) ImportMatch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req ImportMatchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"Invalid request body"}`))
		return
	}

	ctx := r.Context()

	// 1. Resolve or Create Countries
	homeCountryID := h.getOrCreateCountry(ctx, req.HomeTeamName)
	awayCountryID := h.getOrCreateCountry(ctx, req.AwayTeamName)

	// 2. Resolve or Create Teams
	homeTeamID := h.getOrCreateTeam(ctx, homeCountryID, req.HomeTeamName, req.HomeTeamLogoURL)
	awayTeamID := h.getOrCreateTeam(ctx, awayCountryID, req.AwayTeamName, req.AwayTeamLogoURL)

	// 3. Insert Match
	utcTime, err := time.Parse(time.RFC3339, req.UtcDatetime)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"Invalid utc_datetime format (use RFC3339)"}`))
		return
	}

	slug := fmt.Sprintf("wc-2026-%s-vs-%s-%d", makeSlug(req.HomeTeamName), makeSlug(req.AwayTeamName), time.Now().Unix())
	match, err := h.Queries.CreateMatch(ctx, db.CreateMatchParams{
		HomeTeamID:    homeTeamID,
		AwayTeamID:    awayTeamID,
		Title:         req.Title,
		Slug:          slug,
		UtcDatetime:   pgtype.Timestamptz{Time: utcTime, Valid: true},
		Venue:         pgtype.Text{String: req.Venue, Valid: req.Venue != ""},
		CoverImageUrl: pgtype.Text{String: req.CoverImageURL, Valid: req.CoverImageURL != ""},
		HomeScore:     req.HomeScore,
		AwayScore:     req.AwayScore,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(fmt.Sprintf(`{"error":"Internal Error","message":"Failed to create match: %v"}`, err)))
		return
	}

	// Fetch all stat types from database to get their IDs
	statTypes, err := h.Queries.ListStatTypes(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(fmt.Sprintf(`{"error":"Internal Error","message":"Failed to list stat types: %v"}`, err)))
		return
	}
	statTypeMap := make(map[string]pgtype.UUID)
	for _, st := range statTypes {
		statTypeMap[st.Name] = st.ID
	}

	// 4. Insert Performances
	for _, perf := range req.Performances {
		var teamID pgtype.UUID
		var countryID pgtype.UUID
		if perf.TeamName == req.HomeTeamName {
			teamID = homeTeamID
			countryID = homeCountryID
		} else {
			teamID = awayTeamID
			countryID = awayCountryID
		}

		playerID := h.getOrCreatePlayer(ctx, countryID, perf.PlayerName, perf.PhotoURL)
		ptID := h.getOrCreatePlayerTeam(ctx, playerID, teamID, perf.JerseyNumber)

		perfTitle := fmt.Sprintf("%s in %s vs %s", perf.PlayerName, req.HomeTeamName, req.AwayTeamName)
		
		var AverageRating pgtype.Numeric
		_ = AverageRating.Scan(fmt.Sprintf("%.1f", perf.Rating))

		createdPerf, err := h.Queries.CreatePerformance(ctx, db.CreatePerformanceParams{
			MatchID:        match.ID,
			PlayerID:       playerID,
			PlayerTeamID:   ptID,
			Title:          perfTitle,
			JerseyNumber:   pgtype.Int4{Int32: perf.JerseyNumber, Valid: perf.JerseyNumber > 0},
			IsStarter:      true,
			Captain:        false,
			MinutesPlayed:  90,
			AverageRating: AverageRating,
		})
		if err != nil {
			fmt.Printf("Error creating performance for %s: %v\n", perf.PlayerName, err)
			continue
		}

		// List of stats to populate
		statValues := map[string]int32{
			"goals":           perf.Goals,
			"assists":         perf.Assists,
			"saves":           0,
			"tackles":         0,
			"clearances":      0,
			"accurate_passes": 0,
		}

		for statName, val := range statValues {
			if typeID, exists := statTypeMap[statName]; exists {
				var numericVal pgtype.Numeric
				_ = numericVal.Scan(fmt.Sprintf("%d", val))
				_, err = h.Queries.CreatePerformanceStat(ctx, db.CreatePerformanceStatParams{
					PerformanceID: createdPerf.ID,
					StatTypeID:    typeID,
					Value:         numericVal,
				})
				if err != nil {
					fmt.Printf("Error inserting stat %s for player %s: %v\n", statName, perf.PlayerName, err)
				}
			}
		}
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(match)
}

func (h *ImportHandler) getOrCreateCountry(ctx context.Context, name string) pgtype.UUID {
	countries, err := h.Queries.SearchCountries(ctx, name)
	if err == nil && len(countries) > 0 {
		for _, c := range countries {
			if strings.EqualFold(c.Name, name) {
				return c.ID
			}
		}
	}
	
	iso2 := strings.ToUpper(name)
	if len(iso2) > 2 {
		iso2 = iso2[:2]
	}
	iso3 := strings.ToUpper(name)
	if len(iso3) > 3 {
		iso3 = iso3[:3]
	}
	
	newCountry, err := h.Queries.CreateCountry(ctx, db.CreateCountryParams{
		Name: name,
		Iso2: iso2,
		Iso3: iso3,
	})
	if err == nil {
		return newCountry.ID
	}
	return pgtype.UUID{}
}

func (h *ImportHandler) getOrCreateTeam(ctx context.Context, countryID pgtype.UUID, name string, logoURL string) pgtype.UUID {
	slug := makeSlug(name)
	team, err := h.Queries.GetTeamBySlug(ctx, slug)
	if err == nil {
		return team.ID
	}

	shortName := name
	if len(shortName) > 3 {
		shortName = shortName[:3]
	}

	newTeam, err := h.Queries.CreateTeam(ctx, db.CreateTeamParams{
		CountryID: countryID,
		Name:      name,
		ShortName: shortName,
		Slug:      slug,
		Type:      db.TeamTypeNATIONAL,
		LogoUrl:   pgtype.Text{String: logoURL, Valid: logoURL != ""},
	})
	if err == nil {
		return newTeam.ID
	}
	return pgtype.UUID{}
}

func (h *ImportHandler) getOrCreatePlayer(ctx context.Context, countryID pgtype.UUID, name string, photoURL string) pgtype.UUID {
	slug := makeSlug(name) + "-2026"
	player, err := h.Queries.GetPlayerBySlug(ctx, slug)
	if err == nil {
		return player.ID
	}

	newPlayer, err := h.Queries.CreatePlayer(ctx, db.CreatePlayerParams{
		CountryID: countryID,
		Name:      name,
		Slug:      slug,
		KnownAs:   pgtype.Text{String: name, Valid: true},
		PhotoUrl:  pgtype.Text{String: photoURL, Valid: photoURL != ""},
	})
	if err == nil {
		return newPlayer.ID
	}
	return pgtype.UUID{}
}

func (h *ImportHandler) getOrCreatePlayerTeam(ctx context.Context, playerID pgtype.UUID, teamID pgtype.UUID, jersey int32) pgtype.UUID {
	pt, err := h.Queries.JoinPlayerToTeam(ctx, db.JoinPlayerToTeamParams{
		PlayerID:     playerID,
		TeamID:       teamID,
		JerseyNumber: pgtype.Int4{Int32: jersey, Valid: jersey > 0},
		StartDate:    pgtype.Date{Time: time.Now(), Valid: true},
		IsActive:     true,
	})
	if err == nil {
		return pt.ID
	}
	return pgtype.UUID{}
}
