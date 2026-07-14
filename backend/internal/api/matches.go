package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/Nikunjsaini07/performx/backend/internal/db"
)

type MatchesHandler struct {
	Queries *db.Queries
}

type CreateMatchRequest struct {
	HomeTeamID        string  `json:"home_team_id"`
	AwayTeamID        string  `json:"away_team_id"`
	Title             string  `json:"title"`
	Slug              string  `json:"slug"`
	Description       *string `json:"description"`
	Round             *string `json:"round"`
	UtcDatetime       string  `json:"utc_datetime"` // RFC3339
	Venue             *string `json:"venue"`
	CoverImageURL     *string `json:"cover_image_url"`
	HomeScore         *int32  `json:"home_score"`
	AwayScore         *int32  `json:"away_score"`
	HomePenaltyScore  *int32  `json:"home_penalty_score"`
	AwayPenaltyScore  *int32  `json:"away_penalty_score"`
}

type UpdateMatchRequest struct {
	HomeTeamID        *string `json:"home_team_id"`
	AwayTeamID        *string `json:"away_team_id"`
	Title             *string `json:"title"`
	Slug              *string `json:"slug"`
	Description       *string `json:"description"`
	Round             *string `json:"round"`
	UtcDatetime       *string `json:"utc_datetime"`
	Venue             *string `json:"venue"`
	CoverImageURL     *string `json:"cover_image_url"`
	HomeScore         *int32  `json:"home_score"`
	AwayScore         *int32  `json:"away_score"`
	HomePenaltyScore  *int32  `json:"home_penalty_score"`
	AwayPenaltyScore  *int32  `json:"away_penalty_score"`
}

// GET /matches
func (h *MatchesHandler) ListMatches(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	limit, offset := getPaginationParams(r)

	matches, err := h.Queries.ListMatches(r.Context(), db.ListMatchesParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to list matches"}`))
		return
	}

	if matches == nil {
		matches = []db.ListMatchesRow{}
	}

	_ = json.NewEncoder(w).Encode(matches)
}

// GET /matches/search
func (h *MatchesHandler) SearchMatches(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	q := r.URL.Query().Get("q")
	limit, offset := getPaginationParams(r)

	matches, err := h.Queries.SearchMatches(r.Context(), db.SearchMatchesParams{
		Title:  "%" + q + "%",
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to search matches"}`))
		return
	}

	if matches == nil {
		matches = []db.SearchMatchesRow{}
	}

	_ = json.NewEncoder(w).Encode(matches)
}

// GET /matches/{slug}
func (h *MatchesHandler) GetMatch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	slug := r.PathValue("slug")

	match, err := h.Queries.GetMatchBySlug(r.Context(), slug)
	if err != nil {
		if err == pgx.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"error":"Not Found","message":"Match not found"}`))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to fetch match details"}`))
		return
	}

	_ = json.NewEncoder(w).Encode(match)
}

// GET /matches/{slug}/stats
func (h *MatchesHandler) GetMatchStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	slug := r.PathValue("slug")

	match, err := h.Queries.GetMatchBySlug(r.Context(), slug)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"Not Found","message":"Match not found"}`))
		return
	}

	stats, err := h.Queries.GetMatchStats(r.Context(), match.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to fetch match statistics"}`))
		return
	}

	if stats == nil {
		stats = []db.GetMatchStatsRow{}
	}

	_ = json.NewEncoder(w).Encode(stats)
}

// GET /matches/{slug}/performances
func (h *MatchesHandler) GetMatchPerformances(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	slug := r.PathValue("slug")

	match, err := h.Queries.GetMatchBySlug(r.Context(), slug)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"Not Found","message":"Match not found"}`))
		return
	}

	perfs, err := h.Queries.GetMatchPerformances(r.Context(), match.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to fetch player performances"}`))
		return
	}

	if perfs == nil {
		perfs = []db.GetMatchPerformancesRow{}
	}

	_ = json.NewEncoder(w).Encode(perfs)
}

// GET /matches/{slug}/reviews
func (h *MatchesHandler) GetMatchReviews(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	slug := r.PathValue("slug")
	limit, offset := getPaginationParams(r)

	match, err := h.Queries.GetMatchBySlug(r.Context(), slug)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"Not Found","message":"Match not found"}`))
		return
	}

	// Try to get the user ID from context (set by auth middleware if authenticated)
	var viewerID pgtype.UUID
	if userIDStr, ok := r.Context().Value(UserIDKey).(string); ok && userIDStr != "" {
		_ = viewerID.Scan(userIDStr)
	}

	reviews, err := h.Queries.GetMatchReviews(r.Context(), db.GetMatchReviewsParams{
		MatchID:  match.ID,
		Limit:    limit,
		Offset:   offset,
		Column4:  viewerID, // $4 parameter for liked_by_me check
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to fetch reviews"}`))
		return
	}

	if reviews == nil {
		reviews = []db.GetMatchReviewsRow{}
	}

	_ = json.NewEncoder(w).Encode(reviews)
}

// GET /matches/{slug}/ratings
func (h *MatchesHandler) GetMatchRatings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	slug := r.PathValue("slug")
	limit, offset := getPaginationParams(r)

	match, err := h.Queries.GetMatchBySlug(r.Context(), slug)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"Not Found","message":"Match not found"}`))
		return
	}

	summary, err := h.Queries.GetMatchAverageRating(r.Context(), match.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to fetch ratings summary"}`))
		return
	}

	ratings, err := h.Queries.GetMatchRatings(r.Context(), db.GetMatchRatingsParams{
		MatchID: match.ID,
		Limit:   limit,
		Offset:  offset,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to fetch ratings list"}`))
		return
	}

	if ratings == nil {
		ratings = []db.GetMatchRatingsRow{}
	}

	response := map[string]interface{}{
		"average_rating": summary.AverageRating,
		"total_votes":    summary.TotalVotes,
		"ratings":        ratings,
	}

	_ = json.NewEncoder(w).Encode(response)
}

// GET /matches/upcoming
func (h *MatchesHandler) GetUpcomingMatches(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	limit, offset := getPaginationParams(r)

	matches, err := h.Queries.GetUpcomingMatches(r.Context(), db.GetUpcomingMatchesParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to fetch upcoming matches"}`))
		return
	}

	if matches == nil {
		matches = []db.GetUpcomingMatchesRow{}
	}

	_ = json.NewEncoder(w).Encode(matches)
}

// GET /matches/completed
func (h *MatchesHandler) GetCompletedMatches(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	limit, offset := getPaginationParams(r)

	matches, err := h.Queries.GetCompletedMatches(r.Context(), db.GetCompletedMatchesParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to fetch completed matches"}`))
		return
	}

	if matches == nil {
		matches = []db.GetCompletedMatchesRow{}
	}

	_ = json.NewEncoder(w).Encode(matches)
}

// POST /matches
func (h *MatchesHandler) CreateMatch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req CreateMatchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"Invalid request body"}`))
		return
	}

	if req.HomeTeamID == "" || req.AwayTeamID == "" || req.Title == "" || req.Slug == "" || req.UtcDatetime == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"home_team_id, away_team_id, title, slug, and utc_datetime are required"}`))
		return
	}

	var homeTeamID pgtype.UUID
	_ = homeTeamID.Scan(req.HomeTeamID)

	var awayTeamID pgtype.UUID
	_ = awayTeamID.Scan(req.AwayTeamID)

	utcTime, err := time.Parse(time.RFC3339, req.UtcDatetime)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"Invalid utc_datetime format (use RFC3339, e.g. YYYY-MM-DDTHH:MM:SSZ)"}`))
		return
	}

	var homeScore int32
	if req.HomeScore != nil {
		homeScore = *req.HomeScore
	}

	var awayScore int32
	if req.AwayScore != nil {
		awayScore = *req.AwayScore
	}

	match, err := h.Queries.CreateMatch(r.Context(), db.CreateMatchParams{
		HomeTeamID:       homeTeamID,
		AwayTeamID:       awayTeamID,
		Title:            req.Title,
		Slug:             req.Slug,
		Description:      pgtype.Text{String: getString(req.Description), Valid: req.Description != nil},
		Round:            pgtype.Text{String: getString(req.Round), Valid: req.Round != nil},
		UtcDatetime:      pgtype.Timestamptz{Time: utcTime, Valid: true},
		Venue:            pgtype.Text{String: getString(req.Venue), Valid: req.Venue != nil},
		CoverImageUrl:    pgtype.Text{String: getString(req.CoverImageURL), Valid: req.CoverImageURL != nil},
		HomeScore:        homeScore,
		AwayScore:        awayScore,
		HomePenaltyScore: pgtype.Int4{Int32: getInt32(req.HomePenaltyScore), Valid: req.HomePenaltyScore != nil},
		AwayPenaltyScore: pgtype.Int4{Int32: getInt32(req.AwayPenaltyScore), Valid: req.AwayPenaltyScore != nil},
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to create match"}`))
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(match)
}

// PATCH /matches/{slug}
func (h *MatchesHandler) UpdateMatch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	slug := r.PathValue("slug")

	match, err := h.Queries.GetMatchBySlug(r.Context(), slug)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"Not Found","message":"Match to update not found"}`))
		return
	}

	var req UpdateMatchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"Invalid request body"}`))
		return
	}

	homeTeamID := match.HomeTeamID
	if req.HomeTeamID != nil {
		_ = homeTeamID.Scan(*req.HomeTeamID)
	}

	awayTeamID := match.AwayTeamID
	if req.AwayTeamID != nil {
		_ = awayTeamID.Scan(*req.AwayTeamID)
	}

	titleVal := match.Title
	if req.Title != nil {
		titleVal = *req.Title
	}

	slugVal := match.Slug
	if req.Slug != nil {
		slugVal = *req.Slug
	}

	var descriptionVal pgtype.Text
	if req.Description != nil {
		descriptionVal = pgtype.Text{String: *req.Description, Valid: true}
	} else {
		descriptionVal = match.Description
	}

	var roundVal pgtype.Text
	if req.Round != nil {
		roundVal = pgtype.Text{String: *req.Round, Valid: true}
	} else {
		roundVal = match.Round
	}

	var utcTime pgtype.Timestamptz
	if req.UtcDatetime != nil {
		if t, err := time.Parse(time.RFC3339, *req.UtcDatetime); err == nil {
			utcTime = pgtype.Timestamptz{Time: t, Valid: true}
		} else {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"error":"Bad Request","message":"Invalid utc_datetime format"}`))
			return
		}
	} else {
		utcTime = match.UtcDatetime
	}

	var venueVal pgtype.Text
	if req.Venue != nil {
		venueVal = pgtype.Text{String: *req.Venue, Valid: true}
	} else {
		venueVal = match.Venue
	}

	var coverVal pgtype.Text
	if req.CoverImageURL != nil {
		coverVal = pgtype.Text{String: *req.CoverImageURL, Valid: true}
	} else {
		coverVal = match.CoverImageUrl
	}

	homeScoreVal := match.HomeScore
	if req.HomeScore != nil {
		homeScoreVal = *req.HomeScore
	}

	awayScoreVal := match.AwayScore
	if req.AwayScore != nil {
		awayScoreVal = *req.AwayScore
	}

	var homePen pgtype.Int4
	if req.HomePenaltyScore != nil {
		homePen = pgtype.Int4{Int32: *req.HomePenaltyScore, Valid: true}
	} else {
		homePen = match.HomePenaltyScore
	}

	var awayPen pgtype.Int4
	if req.AwayPenaltyScore != nil {
		awayPen = pgtype.Int4{Int32: *req.AwayPenaltyScore, Valid: true}
	} else {
		awayPen = match.AwayPenaltyScore
	}

	updatedMatch, err := h.Queries.UpdateMatch(r.Context(), db.UpdateMatchParams{
		ID:               match.ID,
		HomeTeamID:       homeTeamID,
		AwayTeamID:       awayTeamID,
		Title:            titleVal,
		Slug:             slugVal,
		Description:      descriptionVal,
		Round:            roundVal,
		UtcDatetime:      utcTime,
		Venue:            venueVal,
		CoverImageUrl:    coverVal,
		HomeScore:        homeScoreVal,
		AwayScore:        awayScoreVal,
		HomePenaltyScore: homePen,
		AwayPenaltyScore: awayPen,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to update match"}`))
		return
	}

	_ = json.NewEncoder(w).Encode(updatedMatch)
}

// DELETE /matches/{slug}
func (h *MatchesHandler) DeleteMatch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	slug := r.PathValue("slug")

	match, err := h.Queries.GetMatchBySlug(r.Context(), slug)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"Not Found","message":"Match not found"}`))
		return
	}

	err = h.Queries.DeleteMatch(r.Context(), match.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to delete match"}`))
		return
	}

	_, _ = w.Write([]byte(`{"message":"Match successfully deleted"}`))
}
