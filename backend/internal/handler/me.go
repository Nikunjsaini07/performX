package handler

import (
	"github.com/Nikunjsaini07/performx/backend/internal/middleware"
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/Nikunjsaini07/performx/backend/internal/db"
)

type MeHandler struct {
	Queries *db.Queries
}

type UpdateProfileRequest struct {
	DisplayName *string `json:"display_name,omitempty"`
	Bio         *string `json:"bio,omitempty"`
	AvatarURL   *string `json:"avatar_url,omitempty"`
}

// GET /me
func (h *MeHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract user_id from context
	userIDStr, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":"Unauthorized","message":"User context not found"}`))
		return
	}

	var userID pgtype.UUID
	if err := userID.Scan(userIDStr); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"Invalid user ID format"}`))
		return
	}

	user, err := h.Queries.GetUserByID(r.Context(), userID)
	if err != nil {
		if err == pgx.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"error":"Not Found","message":"User profile not found"}`))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to fetch profile"}`))
		return
	}

	resp := UserPayload{
		ID:          user.ID.String(),
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Email:       user.Email,
		Bio:         user.Bio.String,
		AvatarURL:   user.AvatarUrl.String,
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

// PATCH /me
func (h *MeHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract user_id from context
	userIDStr, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":"Unauthorized","message":"User context not found"}`))
		return
	}

	var userID pgtype.UUID
	if err := userID.Scan(userIDStr); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"Invalid user ID format"}`))
		return
	}

	var req UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"Invalid request body"}`))
		return
	}

	// Prepare nullable updates
	var displayName pgtype.Text
	if req.DisplayName != nil {
		displayName = pgtype.Text{String: *req.DisplayName, Valid: true}
	} else {
		displayName = pgtype.Text{Valid: false}
	}

	var bio pgtype.Text
	if req.Bio != nil {
		bio = pgtype.Text{String: *req.Bio, Valid: true}
	} else {
		bio = pgtype.Text{Valid: false}
	}

	var avatarURL pgtype.Text
	if req.AvatarURL != nil {
		avatarURL = pgtype.Text{String: *req.AvatarURL, Valid: true}
	} else {
		avatarURL = pgtype.Text{Valid: false}
	}

	user, err := h.Queries.UpdateUserProfile(r.Context(), db.UpdateUserProfileParams{
		ID:          userID,
		DisplayName: displayName,
		Bio:         bio,
		AvatarUrl:   avatarURL,
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to update profile"}`))
		return
	}

	resp := UserPayload{
		ID:          user.ID.String(),
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Email:       user.Email,
		Bio:         user.Bio.String,
		AvatarURL:   user.AvatarUrl.String,
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

type UpdateUsernameRequest struct {
	Username string `json:"username"`
}

type UpdateAvatarRequest struct {
	AvatarURL string `json:"avatar_url"`
}

// PATCH /me/username
func (h *MeHandler) UpdateUsername(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userIDStr, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":"Unauthorized","message":"User context not found"}`))
		return
	}

	var userID pgtype.UUID
	_ = userID.Scan(userIDStr)

	var req UpdateUsernameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Username == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"Username cannot be empty"}`))
		return
	}

	// Update username
	user, err := h.Queries.UpdateUsername(r.Context(), db.UpdateUsernameParams{
		ID:       userID,
		Username: req.Username,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to update username (it may already be taken)"}`))
		return
	}

	resp := UserPayload{
		ID:          user.ID.String(),
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Email:       user.Email,
		Bio:         user.Bio.String,
		AvatarURL:   user.AvatarUrl.String,
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

// PATCH /me/avatar
func (h *MeHandler) UpdateAvatar(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userIDStr, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":"Unauthorized","message":"User context not found"}`))
		return
	}

	var userID pgtype.UUID
	_ = userID.Scan(userIDStr)

	var req UpdateAvatarRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"Invalid request body"}`))
		return
	}

	user, err := h.Queries.UpdateAvatar(r.Context(), db.UpdateAvatarParams{
		ID:        userID,
		AvatarUrl: pgtype.Text{String: req.AvatarURL, Valid: true},
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to update avatar"}`))
		return
	}

	resp := UserPayload{
		ID:          user.ID.String(),
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Email:       user.Email,
		Bio:         user.Bio.String,
		AvatarURL:   user.AvatarUrl.String,
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

// DELETE /me
func (h *MeHandler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userIDStr, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":"Unauthorized","message":"User context not found"}`))
		return
	}

	var userID pgtype.UUID
	_ = userID.Scan(userIDStr)

	err := h.Queries.DeleteUser(r.Context(), userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to delete user account"}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"message":"User account successfully deleted"}`))
}
