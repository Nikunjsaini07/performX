package handler

import (
	"github.com/Nikunjsaini07/performx/backend/internal/middleware"
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/Nikunjsaini07/performx/backend/internal/db"
)

type CommentsHandler struct {
	Queries *db.Queries
}

type CommentRequest struct {
	Body string `json:"body"`
}

// POST /match-reviews/{reviewId}/comments
func (h *CommentsHandler) CreateMatchComment(w http.ResponseWriter, r *http.Request) {
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

	var req CommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"Invalid request body"}`))
		return
	}

	if req.Body == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"Comment body is required"}`))
		return
	}

	comment, err := h.Queries.CreateMatchReviewComment(r.Context(), db.CreateMatchReviewCommentParams{
		ReviewID: reviewID,
		UserID:   userID,
		Body:     req.Body,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to submit comment"}`))
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(comment)
}

// GET /match-reviews/{reviewId}/comments
func (h *CommentsHandler) GetMatchComments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	reviewIdStr := r.PathValue("reviewId")
	limit, offset := getPaginationParams(r)

	var reviewID pgtype.UUID
	if err := reviewID.Scan(reviewIdStr); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"Invalid review ID"}`))
		return
	}

	comments, err := h.Queries.GetMatchReviewComments(r.Context(), db.GetMatchReviewCommentsParams{
		ReviewID: reviewID,
		Limit:    limit,
		Offset:   offset,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to fetch comments"}`))
		return
	}

	if comments == nil {
		comments = []db.GetMatchReviewCommentsRow{}
	}

	_ = json.NewEncoder(w).Encode(comments)
}

// PATCH /match-review-comments/{commentId}
func (h *CommentsHandler) UpdateMatchComment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	commentIdStr := r.PathValue("commentId")

	userIDStr, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":"Unauthorized","message":"User not authenticated"}`))
		return
	}

	var commentID pgtype.UUID
	if err := commentID.Scan(commentIdStr); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"Invalid comment ID"}`))
		return
	}

	// Fetch comment first to verify ownership
	comment, err := h.Queries.GetMatchReviewCommentByID(r.Context(), commentID)
	if err != nil {
		if err == pgx.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"error":"Not Found","message":"Comment not found"}`))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to fetch comment details"}`))
		return
	}

	if comment.UserID.String() != userIDStr {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(`{"error":"Forbidden","message":"You are not authorized to edit this comment"}`))
		return
	}

	var req CommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"Invalid request body"}`))
		return
	}

	if req.Body == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"Comment body is required"}`))
		return
	}

	updatedComment, err := h.Queries.UpdateMatchReviewComment(r.Context(), db.UpdateMatchReviewCommentParams{
		ID:   commentID,
		Body: req.Body,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to update comment"}`))
		return
	}

	_ = json.NewEncoder(w).Encode(updatedComment)
}

// DELETE /match-review-comments/{commentId}
func (h *CommentsHandler) DeleteMatchComment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	commentIdStr := r.PathValue("commentId")

	userIDStr, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":"Unauthorized","message":"User not authenticated"}`))
		return
	}

	var commentID pgtype.UUID
	if err := commentID.Scan(commentIdStr); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"Invalid comment ID"}`))
		return
	}

	// Fetch comment first to verify ownership/role
	comment, err := h.Queries.GetMatchReviewCommentByID(r.Context(), commentID)
	if err != nil {
		if err == pgx.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"error":"Not Found","message":"Comment not found"}`))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to fetch comment details"}`))
		return
	}

	// Admins can also delete comments
	var requesterUUID pgtype.UUID
	_ = requesterUUID.Scan(userIDStr)
	requester, _ := h.Queries.GetUserByID(r.Context(), requesterUUID)

	if comment.UserID.String() != userIDStr && requester.Role != db.UserRoleADMIN {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(`{"error":"Forbidden","message":"You are not authorized to delete this comment"}`))
		return
	}

	err = h.Queries.DeleteMatchReviewComment(r.Context(), commentID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to delete comment"}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"message":"Comment successfully deleted"}`))
}

// POST /performance-reviews/{reviewId}/comments
func (h *CommentsHandler) CreatePerformanceComment(w http.ResponseWriter, r *http.Request) {
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

	var req CommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"Invalid request body"}`))
		return
	}

	if req.Body == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"Comment body is required"}`))
		return
	}

	comment, err := h.Queries.CreatePerformanceReviewComment(r.Context(), db.CreatePerformanceReviewCommentParams{
		ReviewID: reviewID,
		UserID:   userID,
		Body:     req.Body,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to submit comment"}`))
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(comment)
}

// GET /performance-reviews/{reviewId}/comments
func (h *CommentsHandler) GetPerformanceComments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	reviewIdStr := r.PathValue("reviewId")
	limit, offset := getPaginationParams(r)

	var reviewID pgtype.UUID
	if err := reviewID.Scan(reviewIdStr); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"Invalid review ID"}`))
		return
	}

	comments, err := h.Queries.GetPerformanceReviewComments(r.Context(), db.GetPerformanceReviewCommentsParams{
		ReviewID: reviewID,
		Limit:    limit,
		Offset:   offset,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to fetch comments"}`))
		return
	}

	if comments == nil {
		comments = []db.GetPerformanceReviewCommentsRow{}
	}

	_ = json.NewEncoder(w).Encode(comments)
}

// PATCH /performance-review-comments/{commentId}
func (h *CommentsHandler) UpdatePerformanceComment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	commentIdStr := r.PathValue("commentId")

	userIDStr, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":"Unauthorized","message":"User not authenticated"}`))
		return
	}

	var commentID pgtype.UUID
	if err := commentID.Scan(commentIdStr); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"Invalid comment ID"}`))
		return
	}

	// Fetch comment first to verify ownership
	comment, err := h.Queries.GetPerformanceReviewCommentByID(r.Context(), commentID)
	if err != nil {
		if err == pgx.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"error":"Not Found","message":"Comment not found"}`))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to fetch comment details"}`))
		return
	}

	if comment.UserID.String() != userIDStr {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(`{"error":"Forbidden","message":"You are not authorized to edit this comment"}`))
		return
	}

	var req CommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"Invalid request body"}`))
		return
	}

	if req.Body == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"Comment body is required"}`))
		return
	}

	updatedComment, err := h.Queries.UpdatePerformanceReviewComment(r.Context(), db.UpdatePerformanceReviewCommentParams{
		ID:   commentID,
		Body: req.Body,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to update comment"}`))
		return
	}

	_ = json.NewEncoder(w).Encode(updatedComment)
}

// DELETE /performance-review-comments/{commentId}
func (h *CommentsHandler) DeletePerformanceComment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	commentIdStr := r.PathValue("commentId")

	userIDStr, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":"Unauthorized","message":"User not authenticated"}`))
		return
	}

	var commentID pgtype.UUID
	if err := commentID.Scan(commentIdStr); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"Invalid comment ID"}`))
		return
	}

	// Fetch comment first to verify ownership/role
	comment, err := h.Queries.GetPerformanceReviewCommentByID(r.Context(), commentID)
	if err != nil {
		if err == pgx.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"error":"Not Found","message":"Comment not found"}`))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to fetch comment details"}`))
		return
	}

	// Admins can also delete comments
	var requesterUUID pgtype.UUID
	_ = requesterUUID.Scan(userIDStr)
	requester, _ := h.Queries.GetUserByID(r.Context(), requesterUUID)

	if comment.UserID.String() != userIDStr && requester.Role != db.UserRoleADMIN {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(`{"error":"Forbidden","message":"You are not authorized to delete this comment"}`))
		return
	}

	err = h.Queries.DeletePerformanceReviewComment(r.Context(), commentID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to delete comment"}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"message":"Comment successfully deleted"}`))
}

// GET /users/{username}/comments
func (h *CommentsHandler) GetUserComments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	username := r.PathValue("username")
	limit, offset := getPaginationParams(r)

	user, err := h.Queries.GetUserByUsername(r.Context(), username)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"Not Found","message":"User not found"}`))
		return
	}

	matchComments, err := h.Queries.GetUserMatchReviewComments(r.Context(), db.GetUserMatchReviewCommentsParams{
		UserID: user.ID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to fetch user's match review comments"}`))
		return
	}

	perfComments, err := h.Queries.GetUserPerformanceReviewComments(r.Context(), db.GetUserPerformanceReviewCommentsParams{
		UserID: user.ID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to fetch user's performance review comments"}`))
		return
	}

	if matchComments == nil {
		matchComments = []db.GetUserMatchReviewCommentsRow{}
	}
	if perfComments == nil {
		perfComments = []db.GetUserPerformanceReviewCommentsRow{}
	}

	resp := map[string]interface{}{
		"match_review_comments":       matchComments,
		"performance_review_comments": perfComments,
	}

	_ = json.NewEncoder(w).Encode(resp)
}
