package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"

	"github.com/jackc/pgx/v5/pgxpool"
)

// genuine-sounding reply snippets used to seed threaded comments on reviews.
var replySnippets = []string{
	"Spot on, couldn't agree more.",
	"Great write-up — that second goal changed everything.",
	"Respectfully I'd have gone a touch lower, but I see your point.",
	"Rewatched it last night and this take aged well.",
	"Best analysis of this one I've read on here.",
	"The pressing detail you flagged is so underrated.",
	"Goosebumps just reading this back.",
	"This is exactly why I signed up. Brilliant.",
	"Hard disagree on the defense but love the passion.",
	"Bookmarking this. Perfectly put.",
}

// seedReviews creates match + performance reviews, ensures each author also has
// a matching rating on that entity, then seeds likes and threaded comments so
// the reviews rank in the trending-reviews feed (which requires likes/comments
// greater than zero). Stored rating columns are finalized by the caller.
func seedReviews(ctx context.Context, pool *pgxpool.Pool, userIDs map[string]string, rng *rand.Rand, dryRun bool) error {
	// ---- Match reviews ----
	for i, mr := range matchReviews {
		actualSlug := resolveMatchSlug(mr.matchSlug)
		var matchID string
		if err := pool.QueryRow(ctx, `SELECT id FROM matches WHERE slug = $1`, actualSlug).Scan(&matchID); err != nil {
			return fmt.Errorf("match review slug %q (%q) not found: %w", mr.matchSlug, actualSlug, err)
		}
		authorID := userIDs[mr.author]

		if dryRun {
			log.Printf("[dry-run] would create match review by %q on %q + likes/comments", mr.author, actualSlug)
			continue
		}

		if _, err := pool.Exec(ctx, `
			INSERT INTO match_ratings (match_id, user_id, rating)
			VALUES ($1, $2, $3::numeric(2,1))
			ON CONFLICT (match_id, user_id) DO UPDATE SET rating = EXCLUDED.rating, updated_at = now()
		`, matchID, authorID, fmt.Sprintf("%.1f", clampRating(mr.rating))); err != nil {
			return fmt.Errorf("author match rating: %w", err)
		}

		var reviewID string
		if err := pool.QueryRow(ctx, `
			INSERT INTO match_reviews (match_id, user_id, title, content)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (match_id, user_id) DO UPDATE SET title = EXCLUDED.title, content = EXCLUDED.content, updated_at = now()
			RETURNING id
		`, matchID, authorID, mr.title, mr.content).Scan(&reviewID); err != nil {
			return fmt.Errorf("insert match review: %w", err)
		}

		likers, commenters := engagers(i, mr.author, rng)
		for _, u := range likers {
			if _, err := pool.Exec(ctx, `
				INSERT INTO match_review_likes (review_id, user_id) VALUES ($1, $2)
				ON CONFLICT (review_id, user_id) DO NOTHING
			`, reviewID, userIDs[u]); err != nil {
				return fmt.Errorf("match review like: %w", err)
			}
		}
		for j, u := range commenters {
			body := replySnippets[(i*3+j)%len(replySnippets)]
			if err := insertMatchComment(ctx, pool, reviewID, userIDs[u], body); err != nil {
				return err
			}
		}
		log.Printf("match review by %q on %q: %d likes, %d comments", mr.author, actualSlug, len(likers), len(commenters))
	}

	// ---- Performance reviews ----
	for i, pr := range perfReviews {
		var perfID string
		if err := pool.QueryRow(ctx, `SELECT id FROM performances WHERE slug = $1`, pr.perfSlug).Scan(&perfID); err != nil {
			return fmt.Errorf("perf review slug %q not found: %w", pr.perfSlug, err)
		}
		authorID := userIDs[pr.author]

		if dryRun {
			log.Printf("[dry-run] would create performance review by %q on %q + likes/comments", pr.author, pr.perfSlug)
			continue
		}

		if _, err := pool.Exec(ctx, `
			INSERT INTO performance_ratings (performance_id, user_id, rating)
			VALUES ($1, $2, $3::numeric(2,1))
			ON CONFLICT (performance_id, user_id) DO UPDATE SET rating = EXCLUDED.rating, updated_at = now()
		`, perfID, authorID, fmt.Sprintf("%.1f", clampRating(pr.rating))); err != nil {
			return fmt.Errorf("author perf rating: %w", err)
		}

		var reviewID string
		if err := pool.QueryRow(ctx, `
			INSERT INTO performance_reviews (performance_id, user_id, title, content)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (performance_id, user_id) DO UPDATE SET title = EXCLUDED.title, content = EXCLUDED.content, updated_at = now()
			RETURNING id
		`, perfID, authorID, pr.title, pr.content).Scan(&reviewID); err != nil {
			return fmt.Errorf("insert perf review: %w", err)
		}

		likers, commenters := engagers(i, pr.author, rng)
		for _, u := range likers {
			if _, err := pool.Exec(ctx, `
				INSERT INTO performance_review_likes (review_id, user_id) VALUES ($1, $2)
				ON CONFLICT (review_id, user_id) DO NOTHING
			`, reviewID, userIDs[u]); err != nil {
				return fmt.Errorf("perf review like: %w", err)
			}
		}
		for j, u := range commenters {
			body := replySnippets[(i*3+j+5)%len(replySnippets)]
			if err := insertPerfComment(ctx, pool, reviewID, userIDs[u], body); err != nil {
				return err
			}
		}
		log.Printf("performance review by %q on %q: %d likes, %d comments", pr.author, pr.perfSlug, len(likers), len(commenters))
	}

	return nil
}

// engagers returns the personas (excluding the author) who like a review and
// the subset who also comment. Counts vary by index so trending review scores
// have a real spread.
func engagers(index int, author string, rng *rand.Rand) (likers []string, commenters []string) {
	likeCount := 3 + (index % 5)    // 3..7 likes
	commentCount := 1 + (index % 3) // 1..3 comments

	var pool []string
	for _, p := range personas {
		if p.username != author {
			pool = append(pool, p.username)
		}
	}
	rng.Shuffle(len(pool), func(a, b int) { pool[a], pool[b] = pool[b], pool[a] })

	if likeCount > len(pool) {
		likeCount = len(pool)
	}
	likers = append(likers, pool[:likeCount]...)

	if commentCount > likeCount {
		commentCount = likeCount
	}
	commenters = append(commenters, likers[:commentCount]...)
	return likers, commenters
}

func insertMatchComment(ctx context.Context, pool *pgxpool.Pool, reviewID, userID, body string) error {
	var exists bool
	if err := pool.QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM match_review_comments WHERE review_id=$1 AND user_id=$2 AND body=$3)
	`, reviewID, userID, body).Scan(&exists); err != nil {
		return err
	}
	if exists {
		return nil
	}
	_, err := pool.Exec(ctx, `
		INSERT INTO match_review_comments (review_id, user_id, body) VALUES ($1, $2, $3)
	`, reviewID, userID, body)
	return err
}

func insertPerfComment(ctx context.Context, pool *pgxpool.Pool, reviewID, userID, body string) error {
	var exists bool
	if err := pool.QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM performance_review_comments WHERE review_id=$1 AND user_id=$2 AND body=$3)
	`, reviewID, userID, body).Scan(&exists); err != nil {
		return err
	}
	if exists {
		return nil
	}
	_, err := pool.Exec(ctx, `
		INSERT INTO performance_review_comments (review_id, user_id, body) VALUES ($1, $2, $3)
	`, reviewID, userID, body)
	return err
}
