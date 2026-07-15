package handler

import (
	"github.com/Nikunjsaini07/performx/backend/internal/middleware"
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/Nikunjsaini07/performx/backend/internal/db"
)

type LikesHandler struct {
	Queries *db.Queries
}

// POST /match-reviews/{reviewId}/like
func (h *LikesHandler) LikeMatchReview(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	reviewIdStr := r.PathValue("reviewId")

	userIDStr, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":"Unauthorized","message":"User not authenticated"}`))
		return
	}

	var userID pgtype.UUID
	_ = userID.Scan(userIDStr)

	var reviewID pgtype.UUID
	if err := reviewID.Scan(reviewIdStr); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"Invalid review ID"}`))
		return
	}

	// Verify review exists
	_, err := h.Queries.GetMatchReviewByID(r.Context(), reviewID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"Not Found","message":"Match review not found"}`))
		return
	}

	// Prevent duplicate likes
	alreadyLiked, err := h.Queries.HasUserLikedMatchReview(r.Context(), db.HasUserLikedMatchReviewParams{
		ReviewID: reviewID,
		UserID:   userID,
	})
	if err == nil && alreadyLiked {
		w.WriteHeader(http.StatusConflict)
		_, _ = w.Write([]byte(`{"error":"Conflict","message":"You have already liked this review"}`))
		return
	}

	like, err := h.Queries.LikeMatchReview(r.Context(), db.LikeMatchReviewParams{
		ReviewID: reviewID,
		UserID:   userID,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to like review"}`))
		return
	}

	_ = json.NewEncoder(w).Encode(like)
}

// DELETE /match-reviews/{reviewId}/like
func (h *LikesHandler) UnlikeMatchReview(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	reviewIdStr := r.PathValue("reviewId")

	userIDStr, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":"Unauthorized","message":"User not authenticated"}`))
		return
	}

	var userID pgtype.UUID
	_ = userID.Scan(userIDStr)

	var reviewID pgtype.UUID
	if err := reviewID.Scan(reviewIdStr); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"Invalid review ID"}`))
		return
	}

	// Verify like exists
	alreadyLiked, err := h.Queries.HasUserLikedMatchReview(r.Context(), db.HasUserLikedMatchReviewParams{
		ReviewID: reviewID,
		UserID:   userID,
	})
	if err != nil || !alreadyLiked {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"Not Found","message":"You have not liked this review"}`))
		return
	}

	err = h.Queries.UnlikeMatchReview(r.Context(), db.UnlikeMatchReviewParams{
		ReviewID: reviewID,
		UserID:   userID,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to unlike review"}`))
		return
	}

	_, _ = w.Write([]byte(`{"message":"Review successfully unliked"}`))
}

// GET /match-reviews/{reviewId}/likes
func (h *LikesHandler) GetMatchReviewLikes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	reviewIdStr := r.PathValue("reviewId")

	var reviewID pgtype.UUID
	if err := reviewID.Scan(reviewIdStr); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"Invalid review ID"}`))
		return
	}

	users, err := h.Queries.GetUsersWhoLikedMatchReview(r.Context(), reviewID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to fetch review likes"}`))
		return
	}

	if users == nil {
		users = []db.GetUsersWhoLikedMatchReviewRow{}
	}

	_ = json.NewEncoder(w).Encode(users)
}

// POST /performance-reviews/{reviewId}/like
func (h *LikesHandler) LikePerformanceReview(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	reviewIdStr := r.PathValue("reviewId")

	userIDStr, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":"Unauthorized","message":"User not authenticated"}`))
		return
	}

	var userID pgtype.UUID
	_ = userID.Scan(userIDStr)

	var reviewID pgtype.UUID
	if err := reviewID.Scan(reviewIdStr); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"Invalid review ID"}`))
		return
	}

	// Verify review exists
	_, err := h.Queries.GetPerformanceReviewByID(r.Context(), reviewID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"Not Found","message":"Performance review not found"}`))
		return
	}

	// Prevent duplicate likes
	alreadyLiked, err := h.Queries.HasUserLikedPerformanceReview(r.Context(), db.HasUserLikedPerformanceReviewParams{
		ReviewID: reviewID,
		UserID:   userID,
	})
	if err == nil && alreadyLiked {
		w.WriteHeader(http.StatusConflict)
		_, _ = w.Write([]byte(`{"error":"Conflict","message":"You have already liked this review"}`))
		return
	}

	like, err := h.Queries.LikePerformanceReview(r.Context(), db.LikePerformanceReviewParams{
		ReviewID: reviewID,
		UserID:   userID,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to like review"}`))
		return
	}

	_ = json.NewEncoder(w).Encode(like)
}

// DELETE /performance-reviews/{reviewId}/like
func (h *LikesHandler) UnlikePerformanceReview(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	reviewIdStr := r.PathValue("reviewId")

	userIDStr, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":"Unauthorized","message":"User not authenticated"}`))
		return
	}

	var userID pgtype.UUID
	_ = userID.Scan(userIDStr)

	var reviewID pgtype.UUID
	if err := reviewID.Scan(reviewIdStr); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"Invalid review ID"}`))
		return
	}

	// Verify like exists
	alreadyLiked, err := h.Queries.HasUserLikedPerformanceReview(r.Context(), db.HasUserLikedPerformanceReviewParams{
		ReviewID: reviewID,
		UserID:   userID,
	})
	if err != nil || !alreadyLiked {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"Not Found","message":"You have not liked this review"}`))
		return
	}

	err = h.Queries.UnlikePerformanceReview(r.Context(), db.UnlikePerformanceReviewParams{
		ReviewID: reviewID,
		UserID:   userID,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to unlike review"}`))
		return
	}

	_, _ = w.Write([]byte(`{"message":"Review successfully unliked"}`))
}

// GET /performance-reviews/{reviewId}/likes
func (h *LikesHandler) GetPerformanceReviewLikes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	reviewIdStr := r.PathValue("reviewId")

	var reviewID pgtype.UUID
	if err := reviewID.Scan(reviewIdStr); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"Invalid review ID"}`))
		return
	}

	users, err := h.Queries.GetUsersWhoLikedPerformanceReview(r.Context(), reviewID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to fetch review likes"}`))
		return
	}

	if users == nil {
		users = []db.GetUsersWhoLikedPerformanceReviewRow{}
	}

	_ = json.NewEncoder(w).Encode(users)
}

// POST /match-review-comments/{commentId}/like
func (h *LikesHandler) LikeMatchComment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	commentIdStr := r.PathValue("commentId")

	userIDStr, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":"Unauthorized","message":"User not authenticated"}`))
		return
	}

	var userID pgtype.UUID
	_ = userID.Scan(userIDStr)

	var commentID pgtype.UUID
	if err := commentID.Scan(commentIdStr); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"Invalid comment ID"}`))
		return
	}

	// Verify comment exists
	_, err := h.Queries.GetMatchReviewCommentByID(r.Context(), commentID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"Not Found","message":"Comment not found"}`))
		return
	}

	// Prevent duplicate likes
	alreadyLiked, err := h.Queries.HasUserLikedMatchReviewComment(r.Context(), db.HasUserLikedMatchReviewCommentParams{
		CommentID: commentID,
		UserID:    userID,
	})
	if err == nil && alreadyLiked {
		w.WriteHeader(http.StatusConflict)
		_, _ = w.Write([]byte(`{"error":"Conflict","message":"You have already liked this comment"}`))
		return
	}

	like, err := h.Queries.LikeMatchReviewComment(r.Context(), db.LikeMatchReviewCommentParams{
		CommentID: commentID,
		UserID:    userID,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to like comment"}`))
		return
	}

	_ = json.NewEncoder(w).Encode(like)
}

// DELETE /match-review-comments/{commentId}/like
func (h *LikesHandler) UnlikeMatchComment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	commentIdStr := r.PathValue("commentId")

	userIDStr, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":"Unauthorized","message":"User not authenticated"}`))
		return
	}

	var userID pgtype.UUID
	_ = userID.Scan(userIDStr)

	var commentID pgtype.UUID
	if err := commentID.Scan(commentIdStr); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"Invalid comment ID"}`))
		return
	}

	// Verify like exists
	alreadyLiked, err := h.Queries.HasUserLikedMatchReviewComment(r.Context(), db.HasUserLikedMatchReviewCommentParams{
		CommentID: commentID,
		UserID:    userID,
	})
	if err != nil || !alreadyLiked {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"Not Found","message":"You have not liked this comment"}`))
		return
	}

	err = h.Queries.UnlikeMatchReviewComment(r.Context(), db.UnlikeMatchReviewCommentParams{
		CommentID: commentID,
		UserID:    userID,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to unlike comment"}`))
		return
	}

	_, _ = w.Write([]byte(`{"message":"Comment successfully unliked"}`))
}

// GET /match-review-comments/{commentId}/likes
func (h *LikesHandler) GetMatchCommentLikes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	commentIdStr := r.PathValue("commentId")

	var commentID pgtype.UUID
	if err := commentID.Scan(commentIdStr); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"Invalid comment ID"}`))
		return
	}

	users, err := h.Queries.GetUsersWhoLikedMatchReviewComment(r.Context(), commentID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to fetch comment likes"}`))
		return
	}

	if users == nil {
		users = []db.GetUsersWhoLikedMatchReviewCommentRow{}
	}

	_ = json.NewEncoder(w).Encode(users)
}

// POST /performance-review-comments/{commentId}/like
func (h *LikesHandler) LikePerformanceComment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	commentIdStr := r.PathValue("commentId")

	userIDStr, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":"Unauthorized","message":"User not authenticated"}`))
		return
	}

	var userID pgtype.UUID
	_ = userID.Scan(userIDStr)

	var commentID pgtype.UUID
	if err := commentID.Scan(commentIdStr); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"Invalid comment ID"}`))
		return
	}

	// Verify comment exists
	_, err := h.Queries.GetPerformanceReviewCommentByID(r.Context(), commentID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"Not Found","message":"Comment not found"}`))
		return
	}

	// Prevent duplicate likes
	alreadyLiked, err := h.Queries.HasUserLikedPerformanceReviewComment(r.Context(), db.HasUserLikedPerformanceReviewCommentParams{
		CommentID: commentID,
		UserID:    userID,
	})
	if err == nil && alreadyLiked {
		w.WriteHeader(http.StatusConflict)
		_, _ = w.Write([]byte(`{"error":"Conflict","message":"You have already liked this comment"}`))
		return
	}

	like, err := h.Queries.LikePerformanceReviewComment(r.Context(), db.LikePerformanceReviewCommentParams{
		CommentID: commentID,
		UserID:    userID,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to like comment"}`))
		return
	}

	_ = json.NewEncoder(w).Encode(like)
}

// DELETE /performance-review-comments/{commentId}/like
func (h *LikesHandler) UnlikePerformanceComment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	commentIdStr := r.PathValue("commentId")

	userIDStr, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":"Unauthorized","message":"User not authenticated"}`))
		return
	}

	var userID pgtype.UUID
	_ = userID.Scan(userIDStr)

	var commentID pgtype.UUID
	if err := commentID.Scan(commentIdStr); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"Invalid comment ID"}`))
		return
	}

	// Verify like exists
	alreadyLiked, err := h.Queries.HasUserLikedPerformanceReviewComment(r.Context(), db.HasUserLikedPerformanceReviewCommentParams{
		CommentID: commentID,
		UserID:    userID,
	})
	if err != nil || !alreadyLiked {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"Not Found","message":"You have not liked this comment"}`))
		return
	}

	err = h.Queries.UnlikePerformanceReviewComment(r.Context(), db.UnlikePerformanceReviewCommentParams{
		CommentID: commentID,
		UserID:    userID,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to unlike comment"}`))
		return
	}

	_, _ = w.Write([]byte(`{"message":"Comment successfully unliked"}`))
}

// GET /performance-review-comments/{commentId}/likes
func (h *LikesHandler) GetPerformanceCommentLikes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	commentIdStr := r.PathValue("commentId")

	var commentID pgtype.UUID
	if err := commentID.Scan(commentIdStr); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"Invalid comment ID"}`))
		return
	}

	users, err := h.Queries.GetUsersWhoLikedPerformanceReviewComment(r.Context(), commentID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to fetch comment likes"}`))
		return
	}

	if users == nil {
		users = []db.GetUsersWhoLikedPerformanceReviewCommentRow{}
	}

	_ = json.NewEncoder(w).Encode(users)
}

