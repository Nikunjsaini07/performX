// Command seed-engagement seeds realistic launch-day community engagement so
// the homepage trending sections (Trending Players, Trending Reviews, Top
// Matches, Top Performances) have genuine content before real users arrive.
//
// It creates 8 realistic user personas (full profiles, avatars, bcrypt
// passwords, verified emails) and has them:
//   - rate marquee matches and the top-rated performances (1-10 scale)
//   - write genuine written reviews on both matches and performances
//   - like and reply to each other's reviews
//
// It then recomputes the stored average_rating/total_votes columns for the
// affected matches/performances and runs one trending computation cycle
// (the exact production worker logic) to populate trending_scores.
//
// The tool is idempotent: users are created only if missing, ratings upsert,
// reviews/likes use ON CONFLICT DO NOTHING, and comments are de-duplicated by
// (review, author, body).
//
// Usage:
//
//	go run ./cmd/seed-engagement -dry-run
//	go run ./cmd/seed-engagement
package main

import (
	"context"
	"flag"
	"log"
	"math/rand"
	"os"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Nikunjsaini07/performx/backend/internal/db"
	"github.com/Nikunjsaini07/performx/backend/internal/worker"
)

type persona struct {
	username    string
	displayName string
	email       string
	bio         string
	avatarURL   string
}

// eight realistic community personas
var personas = []persona{
	{"marco_calcio", "Marco Rossi", "marco.rossi@performx.fans", "Catenaccio believer. I rewatch matches for the shape, not the goals.", "https://i.pravatar.cc/300?u=marco_calcio"},
	{"amara_diallo", "Amara Diallo", "amara.diallo@performx.fans", "Following every African side at the World Cup. AFCON nights raised me.", "https://i.pravatar.cc/300?u=amara_diallo"},
	{"yuki_tanaka", "Yuki Tanaka", "yuki.tanaka@performx.fans", "J-League season ticket holder. Pressing traps are my love language.", "https://i.pravatar.cc/300?u=yuki_tanaka"},
	{"liam_odwyer", "Liam O'Dwyer", "liam.odwyer@performx.fans", "Groundhopper from Cork. 200+ stadiums and counting.", "https://i.pravatar.cc/300?u=liam_odwyer"},
	{"sofia_marquez", "Sofía Márquez", "sofia.marquez@performx.fans", "Argentina till I die. I cried in 2022 and I'll cry again.", "https://i.pravatar.cc/300?u=sofia_marquez"},
	{"hans_becker", "Hans Becker", "hans.becker@performx.fans", "xG spreadsheets and a cold beer. Numbers never lie, players sometimes do.", "https://i.pravatar.cc/300?u=hans_becker"},
	{"priya_nair", "Priya Nair", "priya.nair@performx.fans", "Kerala football kid. Here for the underdogs and the chaos.", "https://i.pravatar.cc/300?u=priya_nair"},
	{"deshawn_carter", "DeShawn Carter", "deshawn.carter@performx.fans", "USMNT diehard. Believe. (Please don't break my heart again.)", "https://i.pravatar.cc/300?u=deshawn_carter"},
}

const seedPassword = "PerformX2026!"

// matchReview is a genuine written match review authored by one persona.
type matchReview struct {
	matchSlug string
	author    string
	rating    float64
	title     string
	content   string
}

// perfReview is a genuine written performance review authored by one persona.
type perfReview struct {
	perfSlug string
	author   string
	rating   float64
	title    string
	content  string
}

var matchReviews = []matchReview{
	{"argentina-vs-switzerland-qf", "sofia_marquez", 9.5, "Extra time drama, vintage Albiceleste",
		"120 minutes of pure heart attack. Switzerland were organised and dangerous on the break, and losing Embolo to that second yellow changed everything. Mac Allister's header set the tone, but it was Álvarez off the bench and Lautaro's stoppage-time dagger that sent us through. This is the kind of night you remember forever."},
	{"norway-vs-england-qf", "hans_becker", 8.5, "Bellingham decides it after the break",
		"Norway pressed brilliantly for an hour and Schjelderup's finish was clinical. But England's midfield control eventually told. Bellingham's brace — one in first-half stoppage, one from the corner in extra time — was the difference. The xG barely separated these teams; quality in the box did."},
	{"spain-vs-belgium-qf", "marco_calcio", 9.0, "A tactical masterclass from Spain",
		"Positional play at its finest. Belgium's golden generation had one last go but Spain suffocated them in midfield and picked the right moments to strike. 2-1 flatters Belgium if anything. The rondo made real."},
	{"argentina-vs-egypt", "amara_diallo", 8.5, "Egypt made them earn every inch",
		"People will remember the Messi comeback but Egypt were magnificent for 60 minutes. Salah caused genuine problems and the Pharaohs went toe to toe. Heartbreaking way to go out 3-2, but they left with heads high."},
	{"usa-vs-belgium", "deshawn_carter", 6.5, "Rough night for the boys",
		"1-4 stings. We competed in patches but Belgium's movement between the lines was too much. De Ketelaere was unplayable. Back to the drawing board, but there's a foundation here. Believe... next time."},
	{"portugal-vs-spain", "marco_calcio", 8.0, "Iberian chess match, tightest of margins",
		"0-1 doesn't capture how tense this was. Two sides who know each other cold, barely a chance until the decisive moment. A game for the purists — I loved every cagey minute of it."},
	{"mexico-vs-england", "liam_odwyer", 8.5, "End to end and worth the ticket",
		"2-3 and it could've been 4-4. Mexico's tempo in the first half was electric, England's quality dragged them back. Neutral's dream. Got to the ground early and stayed till the last whistle — no regrets."},
	{"brazil-vs-norway", "priya_nair", 8.0, "The upset that woke up the group",
		"Norway 2-1 Brazil! Haaland led the line like a man possessed and the underdog spirit was everything I come to this tournament for. Brazil will be fine but tonight belonged to the outsiders."},
}

var perfReviews = []perfReview{
	{"lionel-messi-c675-vs-algeria-763d920d", "sofia_marquez", 9.9, "Simply the maestro",
		"Two assists, a goal, and about six moments where he saw a pass nobody else on the pitch could. 80 minutes of him conducting the whole thing. I don't have the words anymore — just gratitude that we got to watch him."},
	{"jude-bellingham-vs-norway-so5srw", "hans_becker", 9.4, "Big-game player, again",
		"His numbers in knockout football are absurd. Two goals, constant progressive carries, and the maturity to dictate tempo when Norway were pressing. The corner header showed he does the ugly stuff too. Complete midfield performance."},
	{"jonathan-david-5e39-vs-qatar-82ba58b8", "liam_odwyer", 9.5, "Ruthless centre-forward play",
		"Every touch had purpose. David's movement in the box was a class above and he punished Qatar every time they switched off. This is what an elite No.9 looks like."},
	{"charles-de-ketelaere-c8aa-vs-usa-fe4e0dea", "deshawn_carter", 9.0, "He tormented us all night",
		"As a USA fan this hurts to write, but De Ketelaere was sensational. Floated between our lines, created everything, and never gave the ball away. Tip of the cap to a brilliant display."},
	{"vin-cius-j-nior-bce7-vs-scotland-78b8b109", "priya_nair", 9.0, "Unplayable on the left",
		"Scotland doubled up on him and it still didn't matter. Vinícius was pure electricity — beating his man at will and dragging defenders out of shape all game. Box office."},
	{"kylian-mbapp-cb77-vs-sweden-118035bc", "marco_calcio", 9.2, "Devastating in transition",
		"You give Mbappé a yard of space in behind and it's over. Sweden sat deep and he still found the gaps. The acceleration for the decisive moment was frightening. Generational forward."},
}

func main() {
	dryRun := flag.Bool("dry-run", false, "print planned changes without writing to the database")
	force := flag.Bool("force", false, "bypass the already-seeded guard (recovery only; safe because inserts are idempotent and stored ratings are only finalized once)")
	flag.Parse()

	dbURL := loadDBURL()
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set (checked .env and environment)")
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer pool.Close()

	rng := rand.New(rand.NewSource(2026))

	// 1. Users
	userIDs, err := seedUsers(ctx, pool, *dryRun)
	if err != nil {
		log.Fatalf("seedUsers: %v", err)
	}

	// Idempotency guard: if personas already have ratings, don't re-seed
	// (the blended performance-rating math is not safe to run twice).
	if !*dryRun && !*force {
		seeded, err := alreadySeeded(ctx, pool, userIDs)
		if err != nil {
			log.Fatalf("alreadySeeded check: %v", err)
		}
		if seeded {
			log.Println("engagement already seeded (personas have ratings); nothing to do. Use -force to re-run.")
			return
		}
	}

	// 2. Ratings on marquee matches + top performances
	matchIDs, err := seedMatchRatings(ctx, pool, userIDs, distinctMatchSlugs(), rng, *dryRun)
	if err != nil {
		log.Fatalf("seedMatchRatings: %v", err)
	}
	perfSeeds, _, err := seedPerformanceRatings(ctx, pool, userIDs, distinctPerfSlugs(), rng, *dryRun)
	if err != nil {
		log.Fatalf("seedPerformanceRatings: %v", err)
	}

	// 3. Reviews + likes + comments (may add author ratings)
	if err := seedReviews(ctx, pool, userIDs, rng, *dryRun); err != nil {
		log.Fatalf("seedReviews: %v", err)
	}

	if *dryRun {
		log.Println("[dry-run] would finalize stored ratings and run trending computation")
		log.Println("dry run complete, no changes written")
		return
	}

	// 4. Finalize stored average_rating/total_votes now that all ratings
	//    (plan + review authors) are inserted.
	if err := finalizeMatchRatings(ctx, pool, matchIDs); err != nil {
		log.Fatalf("finalizeMatchRatings: %v", err)
	}
	if err := finalizePerformanceRatings(ctx, pool, perfSeeds); err != nil {
		log.Fatalf("finalizePerformanceRatings: %v", err)
	}

	// 5. Recompute trending using the exact production worker logic.
	queries := db.New(pool)
	if err := worker.RunOnce(ctx, pool, queries); err != nil {
		log.Fatalf("trending computation: %v", err)
	}
	log.Println("seed-engagement complete: users, ratings, reviews, likes, comments, and trending scores are in place")
}

func loadDBURL() string {
	if v := os.Getenv("DATABASE_URL"); v != "" {
		return v
	}
	content, err := os.ReadFile(".env")
	if err != nil {
		return ""
	}
	for _, line := range strings.Split(string(content), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		if strings.TrimSpace(parts[0]) == "DATABASE_URL" {
			return strings.Trim(strings.TrimSpace(parts[1]), `'"`)
		}
	}
	return ""
}
