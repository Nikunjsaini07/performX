package handler

import (
	"github.com/Nikunjsaini07/performx/backend/internal/middleware"
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/Nikunjsaini07/performx/backend/internal/db"
)

type ReviewsHandler struct {
	Queries *db.Queries
}

type ReviewRequest struct {
	Title   *string `json:"title"`
	Content string  `json:"content"`
}

type UpdateReviewRequest struct {
	Title   *string `json:"title"`
	Content *string `json:"content"`
}

// POST /matches/{slug}/reviews
func (h *ReviewsHandler) CreateMatchReview(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	slug := r.PathValue("slug")

	userIDStr, ok := r.Context().Value(middleware.UserIDKey).(string)
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

	var req ReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"Invalid request body"}`))
		return
	}

	if req.Content == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"Content is required"}`))
		return
	}

	review, err := h.Queries.CreateMatchReview(r.Context(), db.CreateMatchReviewParams{
		MatchID: match.ID,
		UserID:  userID,
		Title:   pgtype.Text{String: getString(req.Title), Valid: req.Title != nil},
		Content: req.Content,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to submit review"}`))
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(review)
}

// GET /matches/{slug}/reviews/{reviewId}
func (h *ReviewsHandler) GetMatchReview(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	reviewIdStr := r.PathValue("reviewId")

	var reviewID pgtype.UUID
	if err := reviewID.Scan(reviewIdStr); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"Invalid review ID"}`))
		return
	}

	review, err := h.Queries.GetMatchReviewByID(r.Context(), reviewID)
	if err != nil {
		if err == pgx.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"error":"Not Found","message":"Review not found"}`))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to fetch review"}`))
		return
	}

	_ = json.NewEncoder(w).Encode(review)
}

// PATCH /matches/{slug}/reviews/{reviewId}
func (h *ReviewsHandler) UpdateMatchReview(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	reviewIdStr := r.PathValue("reviewId")

	userIDStr, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":"Unauthorized","message":"User not authenticated"}`))
		return
	}

	var reviewID pgtype.UUID
	if err := reviewID.Scan(reviewIdStr); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"Invalid review ID"}`))
		return
	}

	// Fetch review first to verify ownership
	review, err := h.Queries.GetMatchReviewByID(r.Context(), reviewID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"Not Found","message":"Review not found"}`))
		return
	}

	if review.UserID.String() != userIDStr {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(`{"error":"Forbidden","message":"You are not authorized to update this review"}`))
		return
	}

	var req UpdateReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"Invalid request body"}`))
		return
	}

	titleVal := review.Title
	if req.Title != nil {
		titleVal = pgtype.Text{String: *req.Title, Valid: true}
	}

	contentVal := review.Content
	if req.Content != nil {
		contentVal = *req.Content
	}

	updatedReview, err := h.Queries.UpdateMatchReview(r.Context(), db.UpdateMatchReviewParams{
		ID:      reviewID,
		Title:   titleVal,
		Content: contentVal,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to update review"}`))
		return
	}

	_ = json.NewEncoder(w).Encode(updatedReview)
}

// DELETE /matches/{slug}/reviews/{reviewId}
func (h *ReviewsHandler) DeleteMatchReview(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	reviewIdStr := r.PathValue("reviewId")

	userIDStr, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":"Unauthorized","message":"User not authenticated"}`))
		return
	}

	var reviewID pgtype.UUID
	if err := reviewID.Scan(reviewIdStr); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"Invalid review ID"}`))
		return
	}

	// Fetch review first to verify ownership
	review, err := h.Queries.GetMatchReviewByID(r.Context(), reviewID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"Not Found","message":"Review not found"}`))
		return
	}

	// Admins can also delete reviews
	var requesterUUID pgtype.UUID
	_ = requesterUUID.Scan(userIDStr)
	requester, _ := h.Queries.GetUserByID(r.Context(), requesterUUID)

	if review.UserID.String() != userIDStr && requester.Role != db.UserRoleADMIN {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(`{"error":"Forbidden","message":"You are not authorized to delete this review"}`))
		return
	}

	err = h.Queries.DeleteMatchReview(r.Context(), reviewID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to delete review"}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"message":"Review successfully deleted"}`))
}

// POST /performances/{id}/reviews
func (h *ReviewsHandler) CreatePerformanceReview(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := r.PathValue("id")

	userIDStr, ok := r.Context().Value(middleware.UserIDKey).(string)
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
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"Not Found","message":"Performance not found"}`))
		return
	}

	var req ReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"Invalid request body"}`))
		return
	}

	if req.Content == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"Content is required"}`))
		return
	}

	review, err := h.Queries.CreatePerformanceReview(r.Context(), db.CreatePerformanceReviewParams{
		PerformanceID: perfID,
		UserID:        userID,
		Title:         pgtype.Text{String: getString(req.Title), Valid: req.Title != nil},
		Content:       req.Content,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to submit review"}`))
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(review)
}

// GET /performances/{id}/reviews/{reviewId}
func (h *ReviewsHandler) GetPerformanceReview(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	reviewIdStr := r.PathValue("reviewId")

	var reviewID pgtype.UUID
	if err := reviewID.Scan(reviewIdStr); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"Invalid review ID"}`))
		return
	}

	review, err := h.Queries.GetPerformanceReviewByID(r.Context(), reviewID)
	if err != nil {
		if err == pgx.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"error":"Not Found","message":"Review not found"}`))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to fetch review"}`))
		return
	}

	_ = json.NewEncoder(w).Encode(review)
}

// PATCH /performances/{id}/reviews/{reviewId}
func (h *ReviewsHandler) UpdatePerformanceReview(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	reviewIdStr := r.PathValue("reviewId")

	userIDStr, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":"Unauthorized","message":"User not authenticated"}`))
		return
	}

	var reviewID pgtype.UUID
	if err := reviewID.Scan(reviewIdStr); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"Invalid review ID"}`))
		return
	}

	// Fetch review first to verify ownership
	review, err := h.Queries.GetPerformanceReviewByID(r.Context(), reviewID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"Not Found","message":"Review not found"}`))
		return
	}

	if review.UserID.String() != userIDStr {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(`{"error":"Forbidden","message":"You are not authorized to update this review"}`))
		return
	}

	var req UpdateReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"Invalid request body"}`))
		return
	}

	titleVal := review.Title
	if req.Title != nil {
		titleVal = pgtype.Text{String: *req.Title, Valid: true}
	}

	contentVal := review.Content
	if req.Content != nil {
		contentVal = *req.Content
	}

	updatedReview, err := h.Queries.UpdatePerformanceReview(r.Context(), db.UpdatePerformanceReviewParams{
		ID:      reviewID,
		Title:   titleVal,
		Content: contentVal,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to update review"}`))
		return
	}

	_ = json.NewEncoder(w).Encode(updatedReview)
}

// DELETE /performances/{id}/reviews/{reviewId}
func (h *ReviewsHandler) DeletePerformanceReview(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	reviewIdStr := r.PathValue("reviewId")

	userIDStr, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":"Unauthorized","message":"User not authenticated"}`))
		return
	}

	var reviewID pgtype.UUID
	if err := reviewID.Scan(reviewIdStr); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"Invalid review ID"}`))
		return
	}

	// Fetch review first to verify ownership
	review, err := h.Queries.GetPerformanceReviewByID(r.Context(), reviewID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"Not Found","message":"Review not found"}`))
		return
	}

	// Admins can also delete reviews
	var requesterUUID pgtype.UUID
	_ = requesterUUID.Scan(userIDStr)
	requester, _ := h.Queries.GetUserByID(r.Context(), requesterUUID)

	if review.UserID.String() != userIDStr && requester.Role != db.UserRoleADMIN {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(`{"error":"Forbidden","message":"You are not authorized to delete this review"}`))
		return
	}

	err = h.Queries.DeletePerformanceReview(r.Context(), reviewID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to delete review"}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"message":"Review successfully deleted"}`))
}
