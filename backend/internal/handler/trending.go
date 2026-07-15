package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/Nikunjsaini07/performx/backend/internal/db"
)

type TrendingHandler struct {
	Queries *db.Queries
}

type TrendingResponseItem struct {
	EntityID   string      `json:"entity_id"`
	EntityType string      `json:"entity_type"`
	Score      float64     `json:"score"`
	Rank       int32       `json:"rank"`
	TimeWindow string      `json:"time_window"`
	Entity     interface{} `json:"entity,omitempty"`
}

func (h *TrendingHandler) getWindow(r *http.Request) string {
	win := r.URL.Query().Get("window")
	if win == "today" || win == "month" {
		return win
	}
	return "week" // default
}

func (h *TrendingHandler) getLimit(r *http.Request) int32 {
	limit, _ := getPaginationParams(r)
	if limit > 100 {
		return 100
	}
	return limit
}

// -----------------------------------------------------------------
// GET /trending/performances
// -----------------------------------------------------------------
func (h *TrendingHandler) GetTrendingPerformances(w http.ResponseWriter, r *http.Request) {
	h.handleTrending(w, r, "performance", func(ctx context.Context, id pgtype.UUID) (interface{}, error) {
		// Attempt to enrich with performance data
		return h.Queries.GetPerformanceByID(ctx, id)
	})
}

// -----------------------------------------------------------------
// GET /trending/players
// -----------------------------------------------------------------
func (h *TrendingHandler) GetTrendingPlayers(w http.ResponseWriter, r *http.Request) {
	h.handleTrending(w, r, "player", func(ctx context.Context, id pgtype.UUID) (interface{}, error) {
		return h.Queries.GetPlayerByID(ctx, id)
	})
}

// -----------------------------------------------------------------
// GET /trending/matches
// -----------------------------------------------------------------
func (h *TrendingHandler) GetTrendingMatches(w http.ResponseWriter, r *http.Request) {
	h.handleTrending(w, r, "match", func(ctx context.Context, id pgtype.UUID) (interface{}, error) {
		return h.Queries.GetMatchByID(ctx, id)
	})
}

// -----------------------------------------------------------------
// GET /trending/reviews
// -----------------------------------------------------------------
func (h *TrendingHandler) GetTrendingReviews(w http.ResponseWriter, r *http.Request) {
	// For reviews, the ID could be a match_review or performance_review.
	// We'll try match review first, then performance review.
	h.handleTrending(w, r, "review", func(ctx context.Context, id pgtype.UUID) (interface{}, error) {
		mr, err := h.Queries.GetMatchReviewByID(ctx, id)
		if err == nil {
			return mr, nil
		}
		pr, err := h.Queries.GetPerformanceReviewByID(ctx, id)
		if err == nil {
			return pr, nil
		}
		return nil, err
	})
}

// -----------------------------------------------------------------
// Core Handler Logic
// -----------------------------------------------------------------
func (h *TrendingHandler) handleTrending(w http.ResponseWriter, r *http.Request, entityType string, enrichFn func(context.Context, pgtype.UUID) (interface{}, error)) {
	w.Header().Set("Content-Type", "application/json")
	
	window := h.getWindow(r)
	limit := h.getLimit(r)

	scores, err := h.Queries.GetTrendingScores(r.Context(), db.GetTrendingScoresParams{
		EntityType: entityType,
		TimeWindow: window,
		Limit:      limit,
	})
	
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Internal Error", "message": "Failed to fetch trending scores"})
		return
	}

	var response []TrendingResponseItem
	for _, s := range scores {
		item := TrendingResponseItem{
			EntityID:   uuidToString(s.EntityID),
			EntityType: s.EntityType,
			Score:      numericToFloat64(s.Score),
			Rank:       s.Rank,
			TimeWindow: s.TimeWindow,
		}

		if enrichFn != nil {
			if entity, err := enrichFn(r.Context(), s.EntityID); err == nil {
				item.Entity = entity
			}
		}

		response = append(response, item)
	}

	// If no items, return empty array not null
	if response == nil {
		response = []TrendingResponseItem{}
	}

	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"data": response,
		"meta": map[string]interface{}{
			"window": window,
			"limit":  limit,
			"type":   entityType,
		},
	})
}

// Helper to convert UUID to string safely, avoiding import cycles if we define it here
func uuidToString(u pgtype.UUID) string {
	if !u.Valid {
		return ""
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x", u.Bytes[0:4], u.Bytes[4:6], u.Bytes[6:8], u.Bytes[8:10], u.Bytes[10:16])
}

func numericToFloat64(n pgtype.Numeric) float64 {
	if !n.Valid {
		return 0
	}
	f, _ := n.Float64Value()
	return f.Float64
}
