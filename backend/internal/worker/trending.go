package worker

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/Nikunjsaini07/performx/backend/internal/db"
)

type workerCtx struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

var windows = []struct {
	Name     string
	Duration time.Duration
}{
	{"today", 24 * time.Hour},
	{"week", 7 * 24 * time.Hour},
	{"month", 30 * 24 * time.Hour},
}

// RunOnce executes a single trending computation cycle synchronously.
// Useful for seeding, cron jobs, or manual recomputation outside the
// long-running background worker.
func RunOnce(ctx context.Context, pool *pgxpool.Pool, queries *db.Queries) error {
	w := &workerCtx{pool: pool, queries: queries}
	return w.computeAllTrending(ctx)
}

func StartTrendingWorker(pool *pgxpool.Pool, queries *db.Queries, interval time.Duration) {
	w := &workerCtx{pool: pool, queries: queries}

	log.Println("Trending worker: starting background job...")
	
	// Run once immediately on startup in a separate goroutine
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()
		if err := w.computeAllTrending(ctx); err != nil {
			log.Printf("Trending worker: initial computation failed: %v", err)
		}
	}()

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
			if err := w.computeAllTrending(ctx); err != nil {
				log.Printf("Trending worker: computation failed: %v", err)
			}
			cancel()
		}
	}()
}

func (w *workerCtx) computeAllTrending(ctx context.Context) error {
	log.Println("Trending worker: beginning computation cycle")
	startTime := time.Now()

	for _, win := range windows {
		cutoff := startTime.Add(-win.Duration)
		
		if err := w.computePerformances(ctx, win.Name, cutoff); err != nil {
			log.Printf("Trending worker: failed to compute performances for %s: %v", win.Name, err)
		}
		if err := w.computeMatches(ctx, win.Name, cutoff); err != nil {
			log.Printf("Trending worker: failed to compute matches for %s: %v", win.Name, err)
		}
		if err := w.computePlayers(ctx, win.Name, cutoff); err != nil {
			log.Printf("Trending worker: failed to compute players for %s: %v", win.Name, err)
		}
		if err := w.computeReviews(ctx, win.Name, cutoff); err != nil {
			log.Printf("Trending worker: failed to compute reviews for %s: %v", win.Name, err)
		}
		
		log.Printf("Trending worker: computed scores for window=%s", win.Name)
	}

	log.Printf("Trending worker: computation cycle finished in %v", time.Since(startTime))
	return nil
}

// ----------------------------------------------------
// 1. Trending Performances
// ----------------------------------------------------
func (w *workerCtx) computePerformances(ctx context.Context, windowName string, cutoff time.Time) error {
	q := `
		SELECT p.id AS entity_id,
			   COALESCE(SUM(CASE WHEN pr.created_at >= $1 THEN 1 ELSE 0 END), 0) * 2 +
			   COALESCE(SUM(CASE WHEN prev.created_at >= $1 THEN 1 ELSE 0 END), 0) * 10 AS score
		FROM performances p
		LEFT JOIN performance_ratings pr ON pr.performance_id = p.id AND pr.created_at >= $1
		LEFT JOIN performance_reviews prev ON prev.performance_id = p.id AND prev.created_at >= $1
		GROUP BY p.id
		HAVING COUNT(pr.id) + COUNT(prev.id) > 0
		ORDER BY score DESC
		LIMIT 100
	`
	return w.processAggregateQuery(ctx, "performance", windowName, q, cutoff)
}

// ----------------------------------------------------
// 2. Trending Matches
// ----------------------------------------------------
func (w *workerCtx) computeMatches(ctx context.Context, windowName string, cutoff time.Time) error {
	q := `
		SELECT m.id AS entity_id,
			   COUNT(DISTINCT mr.id) * 2 + COUNT(DISTINCT mrev.id) * 10 AS score
		FROM matches m
		LEFT JOIN match_ratings mr ON mr.match_id = m.id AND mr.created_at >= $1
		LEFT JOIN match_reviews mrev ON mrev.match_id = m.id AND mrev.created_at >= $1
		GROUP BY m.id
		HAVING COUNT(mr.id) + COUNT(mrev.id) > 0
		ORDER BY score DESC
		LIMIT 100
	`
	return w.processAggregateQuery(ctx, "match", windowName, q, cutoff)
}

// ----------------------------------------------------
// 3. Trending Players
// ----------------------------------------------------
func (w *workerCtx) computePlayers(ctx context.Context, windowName string, cutoff time.Time) error {
	q := `
		SELECT pl.id AS entity_id,
			   SUM(COALESCE(pr_count.cnt, 0) * 2 + COALESCE(prev_count.cnt, 0) * 10) AS score
		FROM players pl
		JOIN performances p ON p.player_id = pl.id
		LEFT JOIN LATERAL (
			SELECT COUNT(*) AS cnt
			FROM performance_ratings pr
			WHERE pr.performance_id = p.id AND pr.created_at >= $1
		) pr_count ON true
		LEFT JOIN LATERAL (
			SELECT COUNT(*) AS cnt
			FROM performance_reviews prev
			WHERE prev.performance_id = p.id AND prev.created_at >= $1
		) prev_count ON true
		GROUP BY pl.id
		HAVING SUM(COALESCE(pr_count.cnt, 0) + COALESCE(prev_count.cnt, 0)) > 0
		ORDER BY score DESC
		LIMIT 100
	`
	return w.processAggregateQuery(ctx, "player", windowName, q, cutoff)
}

// ----------------------------------------------------
// 4. Trending Reviews
// ----------------------------------------------------
func (w *workerCtx) computeReviews(ctx context.Context, windowName string, cutoff time.Time) error {
	// For reviews, we combine match_reviews and performance_reviews.
	q := `
		WITH combined_reviews AS (
			SELECT mr.id AS entity_id,
				   COUNT(DISTINCT mrl.user_id) * 2 + COUNT(DISTINCT mrc.id) * 5 AS score
			FROM match_reviews mr
			LEFT JOIN match_review_likes mrl ON mrl.review_id = mr.id AND mrl.created_at >= $1
			LEFT JOIN match_review_comments mrc ON mrc.review_id = mr.id AND mrc.created_at >= $1
			GROUP BY mr.id
			HAVING COUNT(mrl.user_id) + COUNT(mrc.id) > 0
			UNION ALL
			SELECT pr.id AS entity_id,
				   COUNT(DISTINCT prl.user_id) * 2 + COUNT(DISTINCT prc.id) * 5 AS score
			FROM performance_reviews pr
			LEFT JOIN performance_review_likes prl ON prl.review_id = pr.id AND prl.created_at >= $1
			LEFT JOIN performance_review_comments prc ON prc.review_id = pr.id AND prc.created_at >= $1
			GROUP BY pr.id
			HAVING COUNT(prl.user_id) + COUNT(prc.id) > 0
		)
		SELECT entity_id, score
		FROM combined_reviews
		ORDER BY score DESC
		LIMIT 100
	`
	return w.processAggregateQuery(ctx, "review", windowName, q, cutoff)
}

// ----------------------------------------------------
// Helper to process raw query, assign ranks, and upsert
// ----------------------------------------------------
func (w *workerCtx) processAggregateQuery(ctx context.Context, entityType, windowName, query string, cutoff time.Time) error {
	rows, err := w.pool.Query(ctx, query, cutoff)
	if err != nil {
		return fmt.Errorf("failed to execute aggregate query: %w", err)
	}
	defer rows.Close()

	type result struct {
		id    pgtype.UUID
		score float64
	}
	
	var results []result
	for rows.Next() {
		var r result
		if err := rows.Scan(&r.id, &r.score); err != nil {
			return fmt.Errorf("failed to scan row: %w", err)
		}
		results = append(results, r)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("rows iteration error: %w", err)
	}

	// 1. Upsert new rankings
	// 2. Keep track of which IDs are currently in the top 100
	validIds := make([]pgtype.UUID, 0, len(results))
	for i, r := range results {
		rank := int32(i + 1)
		validIds = append(validIds, r.id)
		
		// Use pgtype.Numeric for score
		var scoreNumeric pgtype.Numeric
		_ = scoreNumeric.Scan(r.score) // Converts float64 to pgtype.Numeric string internally, but pgx handles it nicely if we just format it as string.

		// Manually map float64 to pgtype.Numeric
		sn := pgtype.Numeric{}
		_ = sn.Scan(fmt.Sprintf("%.2f", r.score))

		err := w.queries.UpsertTrendingScore(ctx, db.UpsertTrendingScoreParams{
			EntityType: entityType,
			EntityID:   r.id,
			TimeWindow: windowName,
			Score:      sn,
			Rank:       rank,
		})
		if err != nil {
			return fmt.Errorf("failed to upsert rank %d: %w", rank, err)
		}
	}

	// 3. Delete stale entries that fell out of the top 100 for this window
	// If validIds is empty, we delete all for this type/window
	if len(validIds) == 0 {
		err := w.queries.DeleteTrendingScoresByTypeAndWindow(ctx, db.DeleteTrendingScoresByTypeAndWindowParams{
			EntityType: entityType,
			TimeWindow: windowName,
		})
		if err != nil {
			return fmt.Errorf("failed to delete all stale entries: %w", err)
		}
		return nil
	}

	// For efficiency, we can just delete where ID NOT IN (validIds). 
	// To do this easily, let's execute a raw query.
	delQuery := `
		DELETE FROM trending_scores
		WHERE entity_type = $1 AND time_window = $2
		  AND entity_id <> ALL($3::uuid[])
	`
	_, err = w.pool.Exec(ctx, delQuery, entityType, windowName, validIds)
	if err != nil {
		return fmt.Errorf("failed to delete stale entries: %w", err)
	}

	return nil
}
