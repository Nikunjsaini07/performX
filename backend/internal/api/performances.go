package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/Nikunjsaini07/performx/backend/internal/db"
)

type PerformancesHandler struct {
	Queries *db.Queries
}

type CreatePerformanceRequest struct {
	MatchID        string   `json:"match_id"`
	PlayerID       string   `json:"player_id"`
	PlayerTeamID   string   `json:"player_team_id"`
	Title          string   `json:"title"`
	Description    *string  `json:"description"`
	CoverImageURL  *string  `json:"cover_image_url"`
	JerseyNumber   *int32   `json:"jersey_number"`
	IsStarter      *bool    `json:"is_starter"`
	Captain        *bool    `json:"captain"`
	MinutesPlayed  *int32   `json:"minutes_played"`
	AverageRating *float64 `json:"average_rating"`
}

type UpdatePerformanceRequest struct {
	MatchID        *string  `json:"match_id"`
	PlayerID       *string  `json:"player_id"`
	PlayerTeamID   *string  `json:"player_team_id"`
	Title          *string  `json:"title"`
	Description    *string  `json:"description"`
	CoverImageURL  *string  `json:"cover_image_url"`
	JerseyNumber   *int32   `json:"jersey_number"`
	IsStarter      *bool    `json:"is_starter"`
	Captain        *bool    `json:"captain"`
	MinutesPlayed  *int32   `json:"minutes_played"`
	AverageRating *float64 `json:"average_rating"`
}

// GET /performances
func (h *PerformancesHandler) ListPerformances(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	limit, offset := getPaginationParams(r)

	perfs, err := h.Queries.ListPerformances(r.Context(), db.ListPerformancesParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to list performances"}`))
		return
	}

	if perfs == nil {
		perfs = []db.ListPerformancesRow{}
	}

	_ = json.NewEncoder(w).Encode(perfs)
}

// GET /performances/search
func (h *PerformancesHandler) SearchPerformances(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	q := r.URL.Query().Get("q")
	limit, offset := getPaginationParams(r)

	perfs, err := h.Queries.SearchPerformances(r.Context(), db.SearchPerformancesParams{
		Title:  "%" + q + "%",
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to search performances"}`))
		return
	}

	if perfs == nil {
		perfs = []db.SearchPerformancesRow{}
	}

	_ = json.NewEncoder(w).Encode(perfs)
}

// GET /performances/{id}
// GetPerformance resolves /performances/{id} by either a UUID or a slug, so
// both id-based and slug-based links work at the same route.
func (h *PerformancesHandler) GetPerformance(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := r.PathValue("id")

	var id pgtype.UUID
	if err := id.Scan(idStr); err == nil {
		perf, err := h.Queries.GetPerformanceByID(r.Context(), id)
		if err == nil {
			_ = json.NewEncoder(w).Encode(perf)
			return
		}
		if err != pgx.ErrNoRows {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to fetch performance details"}`))
			return
		}
		// fall through to not-found below if it was a valid UUID with no match
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"Not Found","message":"Performance record not found"}`))
		return
	}

	// Not a valid UUID — treat the path value as a slug.
	perf, err := h.Queries.GetPerformanceBySlug(r.Context(), pgtype.Text{String: idStr, Valid: true})
	if err != nil {
		if err == pgx.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"error":"Not Found","message":"Performance record not found"}`))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to fetch performance details"}`))
		return
	}

	_ = json.NewEncoder(w).Encode(perf)
}

// GET /performances/{id}/stats
func (h *PerformancesHandler) GetPerformanceStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := r.PathValue("id")

	var id pgtype.UUID
	if err := id.Scan(idStr); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"Invalid performance ID"}`))
		return
	}

	stats, err := h.Queries.GetPerformanceStats(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to fetch stats"}`))
		return
	}

	if stats == nil {
		stats = []db.GetPerformanceStatsRow{}
	}

	_ = json.NewEncoder(w).Encode(stats)
}

// GET /performances/{id}/reviews
func (h *PerformancesHandler) GetPerformanceReviews(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := r.PathValue("id")
	limit, offset := getPaginationParams(r)

	var id pgtype.UUID
	if err := id.Scan(idStr); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"Invalid performance ID"}`))
		return
	}

	reviews, err := h.Queries.GetPerformanceReviews(r.Context(), db.GetPerformanceReviewsParams{
		PerformanceID: id,
		Limit:         limit,
		Offset:        offset,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to fetch reviews"}`))
		return
	}

	if reviews == nil {
		reviews = []db.GetPerformanceReviewsRow{}
	}

	_ = json.NewEncoder(w).Encode(reviews)
}

// GET /performances/{id}/ratings
func (h *PerformancesHandler) GetPerformanceRatings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := r.PathValue("id")
	limit, offset := getPaginationParams(r)

	var id pgtype.UUID
	if err := id.Scan(idStr); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"Invalid performance ID"}`))
		return
	}

	summary, err := h.Queries.GetPerformanceAverageRating(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to fetch ratings summary"}`))
		return
	}

	ratings, err := h.Queries.GetPerformanceRatings(r.Context(), db.GetPerformanceRatingsParams{
		PerformanceID: id,
		Limit:         limit,
		Offset:        offset,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to fetch ratings list"}`))
		return
	}

	if ratings == nil {
		ratings = []db.GetPerformanceRatingsRow{}
	}

	response := map[string]interface{}{
		"average_rating": summary.AverageRating,
		"total_votes":    summary.TotalVotes,
		"ratings":        ratings,
	}

	_ = json.NewEncoder(w).Encode(response)
}

// GET /performances/top-rated
func (h *PerformancesHandler) GetTopRatedPerformances(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	limit, offset := getPaginationParams(r)

	perfs, err := h.Queries.GetTopRatedPerformances(r.Context(), db.GetTopRatedPerformancesParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to fetch top-rated performances"}`))
		return
	}

	if perfs == nil {
		perfs = []db.GetTopRatedPerformancesRow{}
	}

	_ = json.NewEncoder(w).Encode(perfs)
}

// GET /performances/recent
func (h *PerformancesHandler) GetRecentPerformances(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	limit, offset := getPaginationParams(r)

	perfs, err := h.Queries.GetRecentPerformances(r.Context(), db.GetRecentPerformancesParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to fetch recent performances"}`))
		return
	}

	if perfs == nil {
		perfs = []db.GetRecentPerformancesRow{}
	}

	_ = json.NewEncoder(w).Encode(perfs)
}

// POST /performances
func (h *PerformancesHandler) CreatePerformance(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req CreatePerformanceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"Invalid request body"}`))
		return
	}

	if req.MatchID == "" || req.PlayerID == "" || req.PlayerTeamID == "" || req.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"match_id, player_id, player_team_id, and title are required"}`))
		return
	}

	var matchID pgtype.UUID
	_ = matchID.Scan(req.MatchID)

	var playerID pgtype.UUID
	_ = playerID.Scan(req.PlayerID)

	var playerTeamID pgtype.UUID
	_ = playerTeamID.Scan(req.PlayerTeamID)

	var jersey pgtype.Int4
	if req.JerseyNumber != nil {
		jersey = pgtype.Int4{Int32: *req.JerseyNumber, Valid: true}
	}

	var isStarter bool
	if req.IsStarter != nil {
		isStarter = *req.IsStarter
	}

	var captain bool
	if req.Captain != nil {
		captain = *req.Captain
	}

	var minPlayed int32
	if req.MinutesPlayed != nil {
		minPlayed = *req.MinutesPlayed
	}

	var AverageRating pgtype.Numeric
	if req.AverageRating != nil {
		_ = AverageRating.Scan(fmt.Sprintf("%.1f", *req.AverageRating))
	}

	perf, err := h.Queries.CreatePerformance(r.Context(), db.CreatePerformanceParams{
		MatchID:        matchID,
		PlayerID:       playerID,
		PlayerTeamID:   playerTeamID,
		Title:          req.Title,
		Description:    pgtype.Text{String: getString(req.Description), Valid: req.Description != nil},
		CoverImageUrl:  pgtype.Text{String: getString(req.CoverImageURL), Valid: req.CoverImageURL != nil},
		JerseyNumber:   jersey,
		IsStarter:      isStarter,
		Captain:        captain,
		MinutesPlayed:  minPlayed,
		AverageRating: AverageRating,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to create performance record"}`))
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(perf)
}

// PATCH /performances/{id}
func (h *PerformancesHandler) UpdatePerformance(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := r.PathValue("id")

	var id pgtype.UUID
	if err := id.Scan(idStr); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"Invalid performance ID"}`))
		return
	}

	// Fetch existing performance first to resolve COALESCE fallbacks
	perf, err := h.Queries.GetPerformanceByID(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"Not Found","message":"Performance to update not found"}`))
		return
	}

	var req UpdatePerformanceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"Invalid request body"}`))
		return
	}

	matchID := perf.MatchID
	if req.MatchID != nil {
		_ = matchID.Scan(*req.MatchID)
	}

	playerID := perf.PlayerID
	if req.PlayerID != nil {
		_ = playerID.Scan(*req.PlayerID)
	}

	playerTeamID := perf.PlayerTeamID
	if req.PlayerTeamID != nil {
		_ = playerTeamID.Scan(*req.PlayerTeamID)
	}

	titleVal := perf.Title
	if req.Title != nil {
		titleVal = *req.Title
	}

	var descVal pgtype.Text
	if req.Description != nil {
		descVal = pgtype.Text{String: *req.Description, Valid: true}
	} else {
		descVal = perf.Description
	}

	var coverVal pgtype.Text
	if req.CoverImageURL != nil {
		coverVal = pgtype.Text{String: *req.CoverImageURL, Valid: true}
	} else {
		coverVal = perf.CoverImageUrl
	}

	var jerseyVal pgtype.Int4
	if req.JerseyNumber != nil {
		jerseyVal = pgtype.Int4{Int32: *req.JerseyNumber, Valid: true}
	} else {
		jerseyVal = perf.JerseyNumber
	}

	starterVal := perf.IsStarter
	if req.IsStarter != nil {
		starterVal = *req.IsStarter
	}

	captainVal := perf.Captain
	if req.Captain != nil {
		captainVal = *req.Captain
	}

	minVal := perf.MinutesPlayed
	if req.MinutesPlayed != nil {
		minVal = *req.MinutesPlayed
	}

	var ratingVal pgtype.Numeric
	if req.AverageRating != nil {
		_ = ratingVal.Scan(fmt.Sprintf("%.1f", *req.AverageRating))
	} else {
		ratingVal = perf.AverageRating
	}

	updatedPerf, err := h.Queries.UpdatePerformance(r.Context(), db.UpdatePerformanceParams{
		ID:             perf.ID,
		MatchID:        matchID,
		PlayerID:       playerID,
		PlayerTeamID:   playerTeamID,
		Title:          titleVal,
		Description:    descVal,
		CoverImageUrl:  coverVal,
		JerseyNumber:   jerseyVal,
		IsStarter:      starterVal,
		Captain:        captainVal,
		MinutesPlayed:  minVal,
		AverageRating: ratingVal,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to update performance record"}`))
		return
	}

	_ = json.NewEncoder(w).Encode(updatedPerf)
}

// DELETE /performances/{id}
func (h *PerformancesHandler) DeletePerformance(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := r.PathValue("id")

	var id pgtype.UUID
	if err := id.Scan(idStr); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"Invalid performance ID"}`))
		return
	}

	// Fetch first to make sure it exists
	_, err := h.Queries.GetPerformanceByID(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"Not Found","message":"Performance not found"}`))
		return
	}

	err = h.Queries.DeletePerformance(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to delete performance record"}`))
		return
	}

	_, _ = w.Write([]byte(`{"message":"Performance record successfully deleted"}`))
}
