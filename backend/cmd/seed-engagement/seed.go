package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

// seedUsers creates the 8 personas if they don't already exist and returns a
// map of username -> user UUID (as string).
func seedUsers(ctx context.Context, pool *pgxpool.Pool, dryRun bool) (map[string]string, error) {
	ids := make(map[string]string)

	hash, err := bcrypt.GenerateFromPassword([]byte(seedPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	for _, p := range personas {
		var existingID string
		err := pool.QueryRow(ctx, `SELECT id FROM users WHERE username = $1`, p.username).Scan(&existingID)
		if err == nil {
			ids[p.username] = existingID
			log.Printf("user %q already exists (%s)", p.username, existingID)
			continue
		}
		if err != pgx.ErrNoRows {
			return nil, fmt.Errorf("lookup user %q: %w", p.username, err)
		}

		if dryRun {
			log.Printf("[dry-run] would create user %q (%s)", p.username, p.displayName)
			ids[p.username] = "00000000-0000-0000-0000-000000000000"
			continue
		}

		var newID string
		err = pool.QueryRow(ctx, `
			INSERT INTO users (username, display_name, email, bio, avatar_url, password_hash, email_verified, role)
			VALUES ($1, $2, $3, $4, $5, $6, true, 'USER')
			RETURNING id
		`, p.username, p.displayName, p.email, p.bio, p.avatarURL, string(hash)).Scan(&newID)
		if err != nil {
			return nil, fmt.Errorf("insert user %q: %w", p.username, err)
		}
		ids[p.username] = newID
		log.Printf("created user %q (%s)", p.username, newID)
	}

	return ids, nil
}

// alreadySeeded returns true if the personas already have engagement, so a
// re-run doesn't duplicate or drift the blended rating math.
func alreadySeeded(ctx context.Context, pool *pgxpool.Pool, userIDs map[string]string) (bool, error) {
	var cnt int
	err := pool.QueryRow(ctx, `
		SELECT count(*) FROM match_ratings WHERE user_id = ANY($1::uuid[])
	`, userIDValues(userIDs)).Scan(&cnt)
	if err != nil {
		return false, err
	}
	return cnt > 0, nil
}

func userIDValues(userIDs map[string]string) []string {
	out := make([]string, 0, len(userIDs))
	for _, id := range userIDs {
		out = append(out, id)
	}
	return out
}

// ratingPlan deterministically decides which personas rate a given entity and
// with what score, producing varied vote counts so rankings have spread.
func ratingPlan(index int, base float64, rng *rand.Rand) map[string]float64 {
	count := 4 + (index % 5) // 4..8
	if count > len(personas) {
		count = len(personas)
	}
	start := index % len(personas)
	out := make(map[string]float64)
	for i := 0; i < count; i++ {
		p := personas[(start+i)%len(personas)]
		delta := (rng.Float64() - 0.5) * 1.2
		out[p.username] = clampRating(base + delta)
	}
	return out
}

func clampRating(v float64) float64 {
	if v < 1 {
		v = 1
	}
	// The rating column is numeric(2,1): max storable value is 9.9 even though
	// the CHECK constraint permits up to 10.
	if v > 9.9 {
		v = 9.9
	}
	return float64(int(v*10+0.5)) / 10
}

func distinctMatchSlugs() []string {
	return []string{
		"argentina-vs-switzerland-qf",
		"norway-vs-england-qf",
		"spain-vs-belgium-qf",
		"france-vs-morocco-qf",
		"argentina-vs-egypt",
		"usa-vs-belgium",
		"portugal-vs-spain",
		"mexico-vs-england",
		"brazil-vs-norway",
		"paraguay-vs-france",
		"switzerland-vs-colombia",
		"canada-vs-morocco",
	}
}

// matchSlugAliases maps human-friendly slugs to the actual DB slug where they
// differ (older matches use wc-<id> slugs).
var matchSlugAliases = map[string]string{
	"argentina-vs-egypt":      "wc-4653848",
	"usa-vs-belgium":          "wc-4653845",
	"portugal-vs-spain":       "wc-4653844",
	"mexico-vs-england":       "wc-4653847",
	"brazil-vs-norway":        "wc-4653846",
	"paraguay-vs-france":      "wc-4653842",
	"switzerland-vs-colombia": "wc-4653849",
	"canada-vs-morocco":       "wc-4653843",
}

func resolveMatchSlug(s string) string {
	if actual, ok := matchSlugAliases[s]; ok {
		return actual
	}
	return s
}

func distinctPerfSlugs() []string {
	return []string{
		"lionel-messi-c675-vs-algeria-763d920d",
		"jonathan-david-5e39-vs-qatar-82ba58b8",
		"ousmane-demb-l-be95-vs-norway-90b387d7",
		"charles-de-ketelaere-c8aa-vs-usa-fe4e0dea",
		"jude-bellingham-vs-norway-so5srw",
		"kylian-mbapp-cb77-vs-sweden-118035bc",
		"mikel-oyarzabal-a7c3-vs-saudi-arabia-279609dc",
		"ayase-ueda-260e-vs-tunisia-964a9fa6",
		"vin-cius-j-nior-bce7-vs-scotland-78b8b109",
		"eloy-room-a95d-vs-ecuador-355b37de",
		"pape-gueye-c629-vs-iraq-c4ccfafb",
		"cody-gakpo-eaf0-vs-sweden-04c1fe2e",
	}
}

// seedMatchRatings inserts match_ratings for each marquee match. Stored columns
// are finalized separately. Returns the set of touched match IDs.
func seedMatchRatings(ctx context.Context, pool *pgxpool.Pool, userIDs map[string]string, slugs []string, rng *rand.Rand, dryRun bool) ([]string, error) {
	var touched []string
	for i, slug := range slugs {
		actualSlug := resolveMatchSlug(slug)
		var matchID string
		if err := pool.QueryRow(ctx, `SELECT id FROM matches WHERE slug = $1`, actualSlug).Scan(&matchID); err != nil {
			return nil, fmt.Errorf("match slug %q (%q) not found: %w", slug, actualSlug, err)
		}
		touched = append(touched, matchID)

		plan := ratingPlan(i, 8.4, rng)
		if dryRun {
			log.Printf("[dry-run] match %q: would insert %d ratings", actualSlug, len(plan))
			continue
		}
		for username, rating := range plan {
			if _, err := pool.Exec(ctx, `
				INSERT INTO match_ratings (match_id, user_id, rating)
				VALUES ($1, $2, $3::numeric(2,1))
				ON CONFLICT (match_id, user_id) DO UPDATE SET rating = EXCLUDED.rating, updated_at = now()
			`, matchID, userIDs[username], fmt.Sprintf("%.1f", rating)); err != nil {
				return nil, fmt.Errorf("insert match_rating (%s): %w", actualSlug, err)
			}
		}
		log.Printf("match %q: seeded %d ratings", actualSlug, len(plan))
	}
	return touched, nil
}

// seedPerformanceRatings inserts performance_ratings for each top performance.
// It captures each performance's original provider (FotMob) rating BEFORE any
// stored-column mutation, returning perfID -> providerSeed for the finalize
// step. Also returns the ordered list of touched performance IDs.
func seedPerformanceRatings(ctx context.Context, pool *pgxpool.Pool, userIDs map[string]string, slugs []string, rng *rand.Rand, dryRun bool) (map[string]float64, []string, error) {
	seeds := make(map[string]float64)
	var touched []string
	for i, slug := range slugs {
		var perfID string
		var seedRating float64
		if err := pool.QueryRow(ctx, `SELECT id, COALESCE(average_rating, 0)::float8 FROM performances WHERE slug = $1`, slug).Scan(&perfID, &seedRating); err != nil {
			return nil, nil, fmt.Errorf("performance slug %q not found: %w", slug, err)
		}
		seeds[perfID] = seedRating
		touched = append(touched, perfID)

		plan := ratingPlan(i, 8.8, rng)
		if dryRun {
			log.Printf("[dry-run] performance %q: would insert %d ratings (seed=%.1f)", slug, len(plan), seedRating)
			continue
		}
		for username, rating := range plan {
			if _, err := pool.Exec(ctx, `
				INSERT INTO performance_ratings (performance_id, user_id, rating)
				VALUES ($1, $2, $3::numeric(2,1))
				ON CONFLICT (performance_id, user_id) DO UPDATE SET rating = EXCLUDED.rating, updated_at = now()
			`, perfID, userIDs[username], fmt.Sprintf("%.1f", rating)); err != nil {
				return nil, nil, fmt.Errorf("insert performance_rating (%s): %w", slug, err)
			}
		}
		log.Printf("performance %q: seeded %d ratings (seed=%.1f)", slug, len(plan), seedRating)
	}
	return seeds, touched, nil
}

// finalizeMatchRatings sets matches.average_rating = AVG(user ratings) and
// total_votes = COUNT(user ratings). Matches never had a provider seed (the
// 0.0/1 defaults are phantom), so pure user aggregation is correct.
func finalizeMatchRatings(ctx context.Context, pool *pgxpool.Pool, matchIDs []string) error {
	for _, id := range matchIDs {
		if _, err := pool.Exec(ctx, `
			UPDATE matches SET average_rating = sub.avg, total_votes = sub.cnt
			FROM (
				SELECT COALESCE(AVG(rating),0)::numeric(3,1) AS avg, COUNT(*)::int AS cnt
				FROM match_ratings WHERE match_id = $1
			) sub
			WHERE matches.id = $1 AND sub.cnt > 0
		`, id); err != nil {
			return fmt.Errorf("finalize match %s: %w", id, err)
		}
	}
	return nil
}

// finalizePerformanceRatings blends each performance's original provider seed
// (as a single vote) with the seeded user ratings:
//
//	average_rating = (seed + SUM(user_ratings)) / (1 + COUNT(user_ratings))
//	total_votes    = 1 + COUNT(user_ratings)
//
// This honours the design where the provider rating is the initial seed vote
// and community votes accrue on top.
func finalizePerformanceRatings(ctx context.Context, pool *pgxpool.Pool, seeds map[string]float64) error {
	for perfID, seed := range seeds {
		if _, err := pool.Exec(ctx, `
			UPDATE performances SET average_rating = sub.avg, total_votes = sub.cnt
			FROM (
				SELECT
					(($2::numeric + COALESCE(SUM(rating),0)) / (1 + COUNT(*)))::numeric(3,1) AS avg,
					(1 + COUNT(*))::int AS cnt
				FROM performance_ratings WHERE performance_id = $1
			) sub
			WHERE performances.id = $1
		`, perfID, fmt.Sprintf("%.1f", seed)); err != nil {
			return fmt.Errorf("finalize performance %s: %w", perfID, err)
		}
	}
	return nil
}
