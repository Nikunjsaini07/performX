package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/Nikunjsaini07/performx/backend/internal/db"
)

type RatingsHandler struct {
	Queries *db.Queries
}

type RatingRequest struct {
	Rating float64 `json:"rating"`
}

// POST /matches/{slug}/rating
func (h *RatingsHandler) RateMatch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	slug := r.PathValue("slug")

	userIDStr, ok := r.Context().Value(UserIDKey).(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":"Unauthorized","message":"User not authenticated"}`))
		return
	}

	var userID pgtype.UUID
	_ = userID.Scan(userIDStr)

	match, err := h.Queries.GetMatchBySlug(r.Context(), slug)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"Not Found","message":"Match not found"}`))
		return
	}

	var req RatingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"Invalid request body"}`))
		return
	}

	if req.Rating < 0.0 || req.Rating > 10.0 {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"Rating must be between 0.0 and 10.0"}`))
		return
	}

	var ratingNum pgtype.Numeric
	_ = ratingNum.Scan(fmt.Sprintf("%.1f", req.Rating))

	rating, err := h.Queries.RateMatch(r.Context(), db.RateMatchParams{
		MatchID: match.ID,
		UserID:  userID,
		Rating:  ratingNum,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(fmt.Sprintf(`{"error":"Internal Error","message":"Failed to submit rating: %s"}`, err.Error())))
		return
	}

	// Keep the stored average_rating/total_votes columns in sync so the
	// homepage top-lists match the live detail-page value.
	if err := h.Queries.RefreshMatchRating(r.Context(), match.ID); err != nil {
		log.Printf("warning: failed to refresh match rating aggregate: %v", err)
	}

	summary, err := h.Queries.GetMatchAverageRating(r.Context(), match.ID)
	if err != nil {
		log.Printf("warning: failed to fetch updated match rating aggregate: %v", err)
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"rating":         rating.Rating,
		"average_rating": summary.AverageRating,
		"total_votes":    summary.TotalVotes,
	})
}

// PATCH /matches/{slug}/rating
func (h *RatingsHandler) UpdateMatchRating(w http.ResponseWriter, r *http.Request) {
	// Re-use RateMatch as it acts as an upsert (ON CONFLICT DO UPDATE)
	h.RateMatch(w, r)
}

// DELETE /matches/{slug}/rating
func (h *RatingsHandler) DeleteMatchRating(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	slug := r.PathValue("slug")

	userIDStr, ok := r.Context().Value(UserIDKey).(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":"Unauthorized","message":"User not authenticated"}`))
		return
	}

	var userID pgtype.UUID
	_ = userID.Scan(userIDStr)

	match, err := h.Queries.GetMatchBySlug(r.Context(), slug)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"Not Found","message":"Match not found"}`))
		return
	}

	// Verify rating exists first
	_, err = h.Queries.GetUserMatchRating(r.Context(), db.GetUserMatchRatingParams{
		MatchID: match.ID,
		UserID:  userID,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"error":"Not Found","message":"Rating not found"}`))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to fetch rating"}`))
		return
	}

	err = h.Queries.DeleteMatchRating(r.Context(), db.DeleteMatchRatingParams{
		MatchID: match.ID,
		UserID:  userID,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to delete rating"}`))
		return
	}

	if err := h.Queries.RefreshMatchRating(r.Context(), match.ID); err != nil {
		log.Printf("warning: failed to refresh match rating aggregate: %v", err)
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"message":"Rating successfully deleted"}`))
}

// GET /matches/{slug}/ratings/me
func (h *RatingsHandler) GetMyMatchRating(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	slug := r.PathValue("slug")

	userIDStr, ok := r.Context().Value(UserIDKey).(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":"Unauthorized","message":"User not authenticated"}`))
		return
	}

	var userID pgtype.UUID
	_ = userID.Scan(userIDStr)

	match, err := h.Queries.GetMatchBySlug(r.Context(), slug)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"Not Found","message":"Match not found"}`))
		return
	}

	rating, err := h.Queries.GetUserMatchRating(r.Context(), db.GetUserMatchRatingParams{
		MatchID: match.ID,
		UserID:  userID,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"error":"Not Found","message":"You have not rated this match"}`))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to fetch rating"}`))
		return
	}

	_ = json.NewEncoder(w).Encode(rating)
}

// POST /performances/{id}/rating
func (h *RatingsHandler) RatePerformance(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := r.PathValue("id")

	userIDStr, ok := r.Context().Value(UserIDKey).(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":"Unauthorized","message":"User not authenticated"}`))
		return
	}

	var userID pgtype.UUID
	_ = userID.Scan(userIDStr)

	var perfID pgtype.UUID
	if err := perfID.Scan(idStr); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"Invalid performance ID"}`))
		return
	}

	// Verify performance exists
	_, err := h.Queries.GetPerformanceByID(r.Context(), perfID)
	if err != nil {
		if err == pgx.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"error":"Not Found","message":"Performance not found"}`))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to fetch performance details"}`))
		return
	}

	var req RatingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"Invalid request body"}`))
		return
	}

	if req.Rating < 0.0 || req.Rating > 10.0 {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"Rating must be between 0.0 and 10.0"}`))
		return
	}

	var ratingNum pgtype.Numeric
	_ = ratingNum.Scan(fmt.Sprintf("%.1f", req.Rating))

	rating, err := h.Queries.RatePerformance(r.Context(), db.RatePerformanceParams{
		PerformanceID: perfID,
		UserID:        userID,
		Rating:        ratingNum,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(fmt.Sprintf(`{"error":"Internal Error","message":"Failed to submit rating: %s"}`, err.Error())))
		return
	}

	// Keep stored average_rating/total_votes in sync (seed vote + community).
	if err := h.Queries.RefreshPerformanceRating(r.Context(), perfID); err != nil {
		log.Printf("warning: failed to refresh performance rating aggregate: %v", err)
	}

	summary, err := h.Queries.GetPerformanceAverageRating(r.Context(), perfID)
	if err != nil {
		log.Printf("warning: failed to fetch updated performance rating aggregate: %v", err)
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"rating":         rating.Rating,
		"average_rating": summary.AverageRating,
		"total_votes":    summary.TotalVotes,
	})
}

// PATCH /performances/{id}/rating
func (h *RatingsHandler) UpdatePerformanceRating(w http.ResponseWriter, r *http.Request) {
	// Re-use RatePerformance as it acts as an upsert (ON CONFLICT DO UPDATE)
	h.RatePerformance(w, r)
}

// DELETE /performances/{id}/rating
func (h *RatingsHandler) DeletePerformanceRating(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := r.PathValue("id")

	userIDStr, ok := r.Context().Value(UserIDKey).(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":"Unauthorized","message":"User not authenticated"}`))
		return
	}

	var userID pgtype.UUID
	_ = userID.Scan(userIDStr)

	var perfID pgtype.UUID
	if err := perfID.Scan(idStr); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"Invalid performance ID"}`))
		return
	}

	// Verify rating exists first
	_, err := h.Queries.GetUserPerformanceRating(r.Context(), db.GetUserPerformanceRatingParams{
		PerformanceID: perfID,
		UserID:        userID,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"error":"Not Found","message":"Rating not found"}`))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to fetch rating"}`))
		return
	}

	err = h.Queries.DeletePerformanceRating(r.Context(), db.DeletePerformanceRatingParams{
		PerformanceID: perfID,
		UserID:        userID,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to delete rating"}`))
		return
	}

	if err := h.Queries.RefreshPerformanceRating(r.Context(), perfID); err != nil {
		log.Printf("warning: failed to refresh performance rating aggregate: %v", err)
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"message":"Rating successfully deleted"}`))
}

// GET /performances/{id}/ratings/me
func (h *RatingsHandler) GetMyPerformanceRating(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := r.PathValue("id")

	userIDStr, ok := r.Context().Value(UserIDKey).(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":"Unauthorized","message":"User not authenticated"}`))
		return
	}

	var userID pgtype.UUID
	_ = userID.Scan(userIDStr)

	var perfID pgtype.UUID
	if err := perfID.Scan(idStr); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"Invalid performance ID"}`))
		return
	}

	rating, err := h.Queries.GetUserPerformanceRating(r.Context(), db.GetUserPerformanceRatingParams{
		PerformanceID: perfID,
		UserID:        userID,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"error":"Not Found","message":"You have not rated this performance"}`))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to fetch rating"}`))
		return
	}

	_ = json.NewEncoder(w).Encode(rating)
}
