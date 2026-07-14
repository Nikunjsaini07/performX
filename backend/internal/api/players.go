package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/Nikunjsaini07/performx/backend/internal/db"
)

type PlayersHandler struct {
	Queries *db.Queries
}

type CreatePlayerRequest struct {
	Name         string  `json:"name"`
	Slug         string  `json:"slug"`
	FullName     *string `json:"full_name"`
	KnownAs      *string `json:"known_as"`
	DateOfBirth  *string `json:"date_of_birth"` // "YYYY-MM-DD"
	PlaceOfBirth *string `json:"place_of_birth"`
	CountryID    *string `json:"country_id"`
	PhotoURL     *string `json:"photo_url"`
	HeightCM     *int32  `json:"height_cm"`
	WeightKG     *int32  `json:"weight_kg"`
	ShirtName    *string `json:"shirt_name"`
}

type UpdatePlayerRequest struct {
	Name         *string `json:"name"`
	Slug         *string `json:"slug"`
	FullName     *string `json:"full_name"`
	KnownAs      *string `json:"known_as"`
	DateOfBirth  *string `json:"date_of_birth"`
	PlaceOfBirth *string `json:"place_of_birth"`
	CountryID    *string `json:"country_id"`
	PhotoURL     *string `json:"photo_url"`
	HeightCM     *int32  `json:"height_cm"`
	WeightKG     *int32  `json:"weight_kg"`
	ShirtName    *string `json:"shirt_name"`
}

// GET /players
func (h *PlayersHandler) ListPlayers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	limit, offset := getPaginationParams(r)

	players, err := h.Queries.ListPlayers(r.Context(), db.ListPlayersParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to list players"}`))
		return
	}

	if players == nil {
		players = []db.ListPlayersRow{}
	}

	_ = json.NewEncoder(w).Encode(players)
}

// GET /players/search
func (h *PlayersHandler) SearchPlayers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	q := r.URL.Query().Get("q")
	limit, offset := getPaginationParams(r)

	searchTerm := "%"
	if q != "" {
		searchTerm = "%" + q + "%"
	}

	players, err := h.Queries.SearchPlayers(r.Context(), db.SearchPlayersParams{
		Name:   searchTerm,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to search players"}`))
		return
	}

	if players == nil {
		players = []db.SearchPlayersRow{}
	}

	_ = json.NewEncoder(w).Encode(players)
}

// GET /players/{slug}
func (h *PlayersHandler) GetPlayer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	slug := r.PathValue("slug")

	player, err := h.Queries.GetPlayerBySlug(r.Context(), slug)
	if err != nil {
		if err == pgx.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"error":"Not Found","message":"Player not found"}`))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to fetch player"}`))
		return
	}

	_ = json.NewEncoder(w).Encode(player)
}

// GET /players/{slug}/career
func (h *PlayersHandler) GetPlayerCareer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	slug := r.PathValue("slug")

	player, err := h.Queries.GetPlayerBySlug(r.Context(), slug)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"Not Found","message":"Player not found"}`))
		return
	}

	career, err := h.Queries.GetPlayerCareer(r.Context(), player.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to fetch career history"}`))
		return
	}

	if career == nil {
		career = []db.GetPlayerCareerRow{}
	}

	_ = json.NewEncoder(w).Encode(career)
}

// GET /players/{slug}/current-team
func (h *PlayersHandler) GetPlayerCurrentTeam(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	slug := r.PathValue("slug")

	player, err := h.Queries.GetPlayerBySlug(r.Context(), slug)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"Not Found","message":"Player not found"}`))
		return
	}

	currTeam, err := h.Queries.GetCurrentTeam(r.Context(), player.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to fetch active team"}`))
		return
	}

	if currTeam == nil {
		currTeam = []db.GetCurrentTeamRow{}
	}

	_ = json.NewEncoder(w).Encode(currTeam)
}

// GET /players/{slug}/performances
func (h *PlayersHandler) GetPlayerPerformances(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	slug := r.PathValue("slug")
	limit, offset := getPaginationParams(r)

	player, err := h.Queries.GetPlayerBySlug(r.Context(), slug)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"Not Found","message":"Player not found"}`))
		return
	}

	perfs, err := h.Queries.GetPlayerPerformances(r.Context(), db.GetPlayerPerformancesParams{
		PlayerID: player.ID,
		Limit:    limit,
		Offset:   offset,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to fetch performances"}`))
		return
	}

	if perfs == nil {
		perfs = []db.GetPlayerPerformancesRow{}
	}

	_ = json.NewEncoder(w).Encode(perfs)
}

// GET /players/{slug}/stats
func (h *PlayersHandler) GetPlayerStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	slug := r.PathValue("slug")

	player, err := h.Queries.GetPlayerBySlug(r.Context(), slug)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"Not Found","message":"Player not found"}`))
		return
	}

	stats, err := h.Queries.GetPlayerAggregatedStats(r.Context(), player.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to fetch aggregated statistics"}`))
		return
	}

	if stats == nil {
		stats = []db.GetPlayerAggregatedStatsRow{}
	}

	_ = json.NewEncoder(w).Encode(stats)
}

// GET /players/{slug}/reviews
func (h *PlayersHandler) GetPlayerReviews(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	slug := r.PathValue("slug")
	limit, offset := getPaginationParams(r)

	player, err := h.Queries.GetPlayerBySlug(r.Context(), slug)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"Not Found","message":"Player not found"}`))
		return
	}

	reviews, err := h.Queries.GetPlayerReviews(r.Context(), db.GetPlayerReviewsParams{
		PlayerID: player.ID,
		Limit:    limit,
		Offset:   offset,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to fetch reviews"}`))
		return
	}

	if reviews == nil {
		reviews = []db.GetPlayerReviewsRow{}
	}

	_ = json.NewEncoder(w).Encode(reviews)
}

// GET /players/{slug}/ratings
func (h *PlayersHandler) GetPlayerRatings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	slug := r.PathValue("slug")
	limit, offset := getPaginationParams(r)

	player, err := h.Queries.GetPlayerBySlug(r.Context(), slug)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"Not Found","message":"Player not found"}`))
		return
	}

	ratings, err := h.Queries.GetPlayerRatings(r.Context(), db.GetPlayerRatingsParams{
		PlayerID: player.ID,
		Limit:    limit,
		Offset:   offset,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to fetch ratings"}`))
		return
	}

	if ratings == nil {
		ratings = []db.GetPlayerRatingsRow{}
	}

	_ = json.NewEncoder(w).Encode(ratings)
}

// POST /players
func (h *PlayersHandler) CreatePlayer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req CreatePlayerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"Invalid request body"}`))
		return
	}

	if req.Name == "" || req.Slug == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"name, and slug are required"}`))
		return
	}

	var countryID pgtype.UUID
	if req.CountryID != nil {
		_ = countryID.Scan(*req.CountryID)
	}

	var dob pgtype.Date
	if req.DateOfBirth != nil {
		if t, err := time.Parse("2006-01-02", *req.DateOfBirth); err == nil {
			dob = pgtype.Date{Time: t, Valid: true}
		}
	}

	player, err := h.Queries.CreatePlayer(r.Context(), db.CreatePlayerParams{
		Name:         req.Name,
		Slug:         req.Slug,
		FullName:     pgtype.Text{String: getString(req.FullName), Valid: req.FullName != nil},
		KnownAs:      pgtype.Text{String: getString(req.KnownAs), Valid: req.KnownAs != nil},
		DateOfBirth:  dob,
		PlaceOfBirth: pgtype.Text{String: getString(req.PlaceOfBirth), Valid: req.PlaceOfBirth != nil},
		CountryID:    countryID,
		PhotoUrl:     pgtype.Text{String: getString(req.PhotoURL), Valid: req.PhotoURL != nil},
		HeightCm:     pgtype.Int4{Int32: getInt32(req.HeightCM), Valid: req.HeightCM != nil},
		WeightKg:     pgtype.Int4{Int32: getInt32(req.WeightKG), Valid: req.WeightKG != nil},
		ShirtName:    pgtype.Text{String: getString(req.ShirtName), Valid: req.ShirtName != nil},
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to create player record"}`))
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(player)
}

// PATCH /players/{slug}
func (h *PlayersHandler) UpdatePlayer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	slug := r.PathValue("slug")

	player, err := h.Queries.GetPlayerBySlug(r.Context(), slug)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"Not Found","message":"Player to update not found"}`))
		return
	}

	var req UpdatePlayerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"Invalid request body"}`))
		return
	}

	var countryID pgtype.UUID
	if req.CountryID != nil {
		_ = countryID.Scan(*req.CountryID)
	}

	var dob pgtype.Date
	if req.DateOfBirth != nil {
		if t, err := time.Parse("2006-01-02", *req.DateOfBirth); err == nil {
			dob = pgtype.Date{Time: t, Valid: true}
		}
	}

	nameVal := player.Name
	if req.Name != nil {
		nameVal = *req.Name
	}

	slugVal := player.Slug
	if req.Slug != nil {
		slugVal = *req.Slug
	}

	updatedPlayer, err := h.Queries.UpdatePlayer(r.Context(), db.UpdatePlayerParams{
		ID:           player.ID,
		Name:         nameVal,
		Slug:         slugVal,
		FullName:     pgtype.Text{String: getString(req.FullName), Valid: req.FullName != nil},
		KnownAs:      pgtype.Text{String: getString(req.KnownAs), Valid: req.KnownAs != nil},
		DateOfBirth:  dob,
		PlaceOfBirth: pgtype.Text{String: getString(req.PlaceOfBirth), Valid: req.PlaceOfBirth != nil},
		CountryID:    countryID,
		PhotoUrl:     pgtype.Text{String: getString(req.PhotoURL), Valid: req.PhotoURL != nil},
		HeightCm:     pgtype.Int4{Int32: getInt32(req.HeightCM), Valid: req.HeightCM != nil},
		WeightKg:     pgtype.Int4{Int32: getInt32(req.WeightKG), Valid: req.WeightKG != nil},
		ShirtName:    pgtype.Text{String: getString(req.ShirtName), Valid: req.ShirtName != nil},
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to update player record"}`))
		return
	}

	_ = json.NewEncoder(w).Encode(updatedPlayer)
}

// DELETE /players/{slug}
func (h *PlayersHandler) DeletePlayer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	slug := r.PathValue("slug")

	player, err := h.Queries.GetPlayerBySlug(r.Context(), slug)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"Not Found","message":"Player not found"}`))
		return
	}

	err = h.Queries.DeletePlayer(r.Context(), player.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to delete player"}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"message":"Player successfully deleted"}`))
}

// Struct pointer extraction helpers
func getString(p *string) string {
	if p != nil {
		return *p
	}
	return ""
}

func getInt32(p *int32) int32 {
	if p != nil {
		return *p
	}
	return 0
}
