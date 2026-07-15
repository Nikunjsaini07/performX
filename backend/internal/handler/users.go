package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/Nikunjsaini07/performx/backend/internal/db"
)

type UsersHandler struct {
	Queries   *db.Queries
	JWTSecret []byte
}

type ProfileResponse struct {
	ID                 string `json:"id"`
	Username           string `json:"username"`
	DisplayName        string `json:"display_name"`
	Email              string `json:"email"`
	Bio                string `json:"bio"`
	AvatarURL          string `json:"avatar_url"`
	CreatedAt          string `json:"created_at"`
	ReviewCount        int64  `json:"review_count"`
	RatingCount        int64  `json:"rating_count"`
	LikesReceivedCount int64  `json:"likes_received_count"`
}

// Helper to parse pagination params
func getPaginationParams(r *http.Request) (int32, int32) {
	limit := int32(20)
	offset := int32(0)

	if lVal := r.URL.Query().Get("limit"); lVal != "" {
		if l, err := strconv.Atoi(lVal); err == nil && l > 0 {
			limit = int32(l)
		}
	}
	if oVal := r.URL.Query().Get("offset"); oVal != "" {
		if o, err := strconv.Atoi(oVal); err == nil && o >= 0 {
			offset = int32(o)
		}
	}
	return limit, offset
}

// GET /users
func (h *UsersHandler) SearchUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	q := r.URL.Query().Get("q")
	limit, offset := getPaginationParams(r)

	searchTerm := "%"
	if q != "" {
		searchTerm = "%" + q + "%"
	}

	users, err := h.Queries.SearchUsers(r.Context(), db.SearchUsersParams{
		Username: searchTerm,
		Limit:    limit,
		Offset:   offset,
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to query users"}`))
		return
	}

	if users == nil {
		users = []db.SearchUsersRow{}
	}

	_ = json.NewEncoder(w).Encode(users)
}

// GET /users/{username}
func (h *UsersHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	username := r.PathValue("username")
	if username == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"Username parameter is required"}`))
		return
	}

	user, err := h.Queries.GetUserByUsername(r.Context(), username)
	if err != nil {
		if err == pgx.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"error":"Not Found","message":"User profile not found"}`))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to fetch user"}`))
		return
	}

	// Fetch Stats
	reviewCount, _ := h.Queries.GetUserReviewCount(r.Context(), user.ID)
	ratingCount, _ := h.Queries.GetUserRatingCount(r.Context(), user.ID)
	likesCount, _ := h.Queries.GetUserLikesReceived(r.Context(), user.ID)

	resp := ProfileResponse{
		ID:                 user.ID.String(),
		Username:           user.Username,
		DisplayName:        user.DisplayName,
		Email:              user.Email,
		Bio:                user.Bio.String,
		AvatarURL:          user.AvatarUrl.String,
		CreatedAt:          user.CreatedAt.Time.Format(time.RFC3339),
		ReviewCount:        reviewCount,
		RatingCount:        ratingCount,
		LikesReceivedCount: likesCount,
	}

	_ = json.NewEncoder(w).Encode(resp)
}

// GET /users/{username}/reviews
func (h *UsersHandler) GetReviews(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	username := r.PathValue("username")
	limit, offset := getPaginationParams(r)

	user, err := h.Queries.GetUserByUsername(r.Context(), username)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"Not Found","message":"User not found"}`))
		return
	}

	reviews, err := h.Queries.GetUserRecentReviews(r.Context(), db.GetUserRecentReviewsParams{
		UserID: user.ID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to fetch reviews"}`))
		return
	}

	if reviews == nil {
		reviews = []db.GetUserRecentReviewsRow{}
	}

	_ = json.NewEncoder(w).Encode(reviews)
}

// GET /users/{username}/ratings
func (h *UsersHandler) GetRatings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	username := r.PathValue("username")
	limit, offset := getPaginationParams(r)

	user, err := h.Queries.GetUserByUsername(r.Context(), username)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"Not Found","message":"User not found"}`))
		return
	}

	ratings, err := h.Queries.GetUserRecentRatings(r.Context(), db.GetUserRecentRatingsParams{
		UserID: user.ID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to fetch ratings"}`))
		return
	}

	if ratings == nil {
		ratings = []db.GetUserRecentRatingsRow{}
	}

	_ = json.NewEncoder(w).Encode(ratings)
}

// GET /users/{username}/activity
func (h *UsersHandler) GetActivity(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	username := r.PathValue("username")
	limit, offset := getPaginationParams(r)

	user, err := h.Queries.GetUserByUsername(r.Context(), username)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"Not Found","message":"User not found"}`))
		return
	}

	activity, err := h.Queries.GetUserActivity(r.Context(), db.GetUserActivityParams{
		UserID: user.ID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(fmt.Sprintf(`{"error":"Internal Error","message":"Failed to fetch activity: %s"}`, err.Error())))
		return
	}

	if activity == nil {
		activity = []db.GetUserActivityRow{}
	}

	_ = json.NewEncoder(w).Encode(activity)
}
