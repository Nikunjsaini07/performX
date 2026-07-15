package handler

import (
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/Nikunjsaini07/performx/backend/internal/db"
)

type TeamsHandler struct {
	Queries *db.Queries
}

type CreateTeamRequest struct {
	CountryID   string  `json:"country_id"`
	Name        string  `json:"name"`
	ShortName   string  `json:"short_name"`
	Slug        string  `json:"slug"`
	Type        string  `json:"type"` // "CLUB" or "NATIONAL"
	LogoURL     *string `json:"logo_url"`
	FoundedYear *int32  `json:"founded_year"`
}

type UpdateTeamRequest struct {
	CountryID   *string `json:"country_id"`
	Name        *string `json:"name"`
	ShortName   *string `json:"short_name"`
	Slug        *string `json:"slug"`
	Type        *string `json:"type"`
	LogoURL     *string `json:"logo_url"`
	FoundedYear *int32  `json:"founded_year"`
}

// GET /teams
func (h *TeamsHandler) ListTeams(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	limit, offset := getPaginationParams(r)

	teams, err := h.Queries.ListTeams(r.Context(), db.ListTeamsParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to list teams"}`))
		return
	}

	if teams == nil {
		teams = []db.ListTeamsRow{}
	}

	_ = json.NewEncoder(w).Encode(teams)
}

// GET /teams/search
func (h *TeamsHandler) SearchTeams(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	q := r.URL.Query().Get("q")
	limit, offset := getPaginationParams(r)

	searchTerm := "%"
	if q != "" {
		searchTerm = "%" + q + "%"
	}

	teams, err := h.Queries.SearchTeams(r.Context(), db.SearchTeamsParams{
		Name:   searchTerm,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to search teams"}`))
		return
	}

	if teams == nil {
		teams = []db.SearchTeamsRow{}
	}

	_ = json.NewEncoder(w).Encode(teams)
}

// GET /teams/{slug}
func (h *TeamsHandler) GetTeam(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	slug := r.PathValue("slug")

	team, err := h.Queries.GetTeamBySlug(r.Context(), slug)
	if err != nil {
		if err == pgx.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"error":"Not Found","message":"Team not found"}`))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to fetch team"}`))
		return
	}

	_ = json.NewEncoder(w).Encode(team)
}

// GET /teams/{slug}/players
func (h *TeamsHandler) GetTeamPlayers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	slug := r.PathValue("slug")

	team, err := h.Queries.GetTeamBySlug(r.Context(), slug)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"Not Found","message":"Team not found"}`))
		return
	}

	players, err := h.Queries.GetTeamPlayers(r.Context(), team.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to fetch team players"}`))
		return
	}

	if players == nil {
		players = []db.GetTeamPlayersRow{}
	}

	_ = json.NewEncoder(w).Encode(players)
}

// GET /teams/{slug}/matches
func (h *TeamsHandler) GetTeamMatches(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	slug := r.PathValue("slug")
	limit, offset := getPaginationParams(r)

	team, err := h.Queries.GetTeamBySlug(r.Context(), slug)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"Not Found","message":"Team not found"}`))
		return
	}

	matches, err := h.Queries.GetTeamMatches(r.Context(), db.GetTeamMatchesParams{
		HomeTeamID: team.ID,
		Limit:      limit,
		Offset:     offset,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to fetch matches"}`))
		return
	}

	if matches == nil {
		matches = []db.GetTeamMatchesRow{}
	}

	_ = json.NewEncoder(w).Encode(matches)
}

// GET /teams/{slug}/performances
func (h *TeamsHandler) GetTeamPerformances(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	slug := r.PathValue("slug")
	limit, offset := getPaginationParams(r)

	team, err := h.Queries.GetTeamBySlug(r.Context(), slug)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"Not Found","message":"Team not found"}`))
		return
	}

	perfs, err := h.Queries.GetTeamPerformances(r.Context(), db.GetTeamPerformancesParams{
		TeamID: team.ID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to fetch performances"}`))
		return
	}

	if perfs == nil {
		perfs = []db.GetTeamPerformancesRow{}
	}

	_ = json.NewEncoder(w).Encode(perfs)
}

// GET /teams/{slug}/stats
func (h *TeamsHandler) GetTeamStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	slug := r.PathValue("slug")

	team, err := h.Queries.GetTeamBySlug(r.Context(), slug)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"Not Found","message":"Team not found"}`))
		return
	}

	stats, err := h.Queries.GetTeamAggregatedStats(r.Context(), team.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to fetch aggregated team statistics"}`))
		return
	}

	if stats == nil {
		stats = []db.GetTeamAggregatedStatsRow{}
	}

	_ = json.NewEncoder(w).Encode(stats)
}

// GET /teams/{slug}/reviews
func (h *TeamsHandler) GetTeamReviews(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	slug := r.PathValue("slug")
	limit, offset := getPaginationParams(r)

	team, err := h.Queries.GetTeamBySlug(r.Context(), slug)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"Not Found","message":"Team not found"}`))
		return
	}

	reviews, err := h.Queries.GetTeamReviews(r.Context(), db.GetTeamReviewsParams{
		TeamID: team.ID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to fetch reviews"}`))
		return
	}

	if reviews == nil {
		reviews = []db.GetTeamReviewsRow{}
	}

	_ = json.NewEncoder(w).Encode(reviews)
}

// GET /teams/{slug}/ratings
func (h *TeamsHandler) GetTeamRatings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	slug := r.PathValue("slug")
	limit, offset := getPaginationParams(r)

	team, err := h.Queries.GetTeamBySlug(r.Context(), slug)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"Not Found","message":"Team not found"}`))
		return
	}

	ratings, err := h.Queries.GetTeamRatings(r.Context(), db.GetTeamRatingsParams{
		HomeTeamID: team.ID,
		Limit:      limit,
		Offset:     offset,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to fetch ratings"}`))
		return
	}

	if ratings == nil {
		ratings = []db.GetTeamRatingsRow{}
	}

	_ = json.NewEncoder(w).Encode(ratings)
}

// POST /teams
func (h *TeamsHandler) CreateTeam(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req CreateTeamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"Invalid request body"}`))
		return
	}

	if req.CountryID == "" || req.Name == "" || req.ShortName == "" || req.Slug == "" || req.Type == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"country_id, name, short_name, slug, and type are required"}`))
		return
	}

	var countryID pgtype.UUID
	_ = countryID.Scan(req.CountryID)

	team, err := h.Queries.CreateTeam(r.Context(), db.CreateTeamParams{
		CountryID:   countryID,
		Name:        req.Name,
		ShortName:   req.ShortName,
		Slug:        req.Slug,
		Type:        db.TeamType(req.Type),
		LogoUrl:     pgtype.Text{String: getString(req.LogoURL), Valid: req.LogoURL != nil},
		FoundedYear: pgtype.Int4{Int32: getInt32(req.FoundedYear), Valid: req.FoundedYear != nil},
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to create team record"}`))
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(team)
}

// PATCH /teams/{slug}
func (h *TeamsHandler) UpdateTeam(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	slug := r.PathValue("slug")

	team, err := h.Queries.GetTeamBySlug(r.Context(), slug)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"Not Found","message":"Team to update not found"}`))
		return
	}

	var req UpdateTeamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"Invalid request body"}`))
		return
	}

	// Update only changed values
	countryIDVal := team.CountryID
	if req.CountryID != nil {
		_ = countryIDVal.Scan(*req.CountryID)
	}

	nameVal := team.Name
	if req.Name != nil {
		nameVal = *req.Name
	}

	shortNameVal := team.ShortName
	if req.ShortName != nil {
		shortNameVal = *req.ShortName
	}

	slugVal := team.Slug
	if req.Slug != nil {
		slugVal = *req.Slug
	}

	typeVal := team.Type
	if req.Type != nil {
		typeVal = db.TeamType(*req.Type)
	}

	updatedTeam, err := h.Queries.UpdateTeam(r.Context(), db.UpdateTeamParams{
		ID:          team.ID,
		CountryID:   countryIDVal,
		Name:        nameVal,
		ShortName:   shortNameVal,
		Slug:        slugVal,
		Type:        typeVal,
		LogoUrl:     pgtype.Text{String: getString(req.LogoURL), Valid: req.LogoURL != nil},
		FoundedYear: pgtype.Int4{Int32: getInt32(req.FoundedYear), Valid: req.FoundedYear != nil},
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to update team record"}`))
		return
	}

	_ = json.NewEncoder(w).Encode(updatedTeam)
}

// DELETE /teams/{slug}
func (h *TeamsHandler) DeleteTeam(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	slug := r.PathValue("slug")

	team, err := h.Queries.GetTeamBySlug(r.Context(), slug)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"Not Found","message":"Team not found"}`))
		return
	}

	err = h.Queries.DeleteTeam(r.Context(), team.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to delete team"}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"message":"Team successfully deleted"}`))
}
