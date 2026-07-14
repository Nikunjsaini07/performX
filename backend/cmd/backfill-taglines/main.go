// Command backfill-taglines reads all matches and performances from the DB,
// uses Serper.dev for Google search context, and Groq (Llama 3) to generate
// high-quality, contextual taglines and descriptions, then writes them back.
//
// Usage:
//
//	go run ./cmd/backfill-taglines -dry-run    # preview without writing
//	go run ./cmd/backfill-taglines             # apply changes
//	go run ./cmd/backfill-taglines -matches    # only update matches
//	go run ./cmd/backfill-taglines -perfs      # only update performances
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

// ── Serper API types ────────────────────────────────────────────────────────

type serperRequest struct {
	Q string `json:"q"`
}

type serperResponse struct {
	Organic []struct {
		Title   string `json:"title"`
		Snippet string `json:"snippet"`
	} `json:"organic"`
}

// ── Groq API types ────────────────────────────────────────────────────────────

type groqRequest struct {
	Model       string        `json:"model"`
	Messages    []groqMessage `json:"messages"`
	Temperature float64       `json:"temperature"`
	MaxTokens   int           `json:"max_tokens"`
}

type groqMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type groqResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

// ── DB row types ────────────────────────────────────────────────────────────

type matchRow struct {
	ID           string
	Title        string
	HomeTeamName string
	AwayTeamName string
	HomeScore    int
	AwayScore    int
	HomePenalty  *int
	AwayPenalty  *int
	Round        *string
	Venue        *string
	UtcDatetime  time.Time
}

type perfRow struct {
	ID            string
	PlayerName    string
	TeamName      string
	MatchTitle    string
	HomeTeamName  string
	AwayTeamName  string
	HomeScore     int
	AwayScore     int
	Round         *string
	Rating        *float64
	Goals         int
	Assists       int
	MinutesPlayed int
	IsStarter     bool
	Captain       bool
	JerseyNumber  *int
}

// ── Generated result ────────────────────────────────────────────────────────

type generatedItem struct {
	ID          string `json:"id"`
	Tagline     string `json:"tagline"`
	Description string `json:"description"`
}

// ── Main ────────────────────────────────────────────────────────────────────

func main() {
	dryRun := flag.Bool("dry-run", false, "preview generated taglines without writing to DB")
	matchesOnly := flag.Bool("matches", false, "only update matches")
	perfsOnly := flag.Bool("perfs", false, "only update performances")
	flag.Parse()

	if err := godotenv.Load(".env"); err != nil {
		log.Println("no .env file found, relying on system env")
	}

	dbURL := os.Getenv("DATABASE_URL")
	groqKey := os.Getenv("GROQ_API_KEY")
	serperKey := os.Getenv("SERPER_API_KEY")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}
	if groqKey == "" {
		log.Fatal("GROQ_API_KEY is not set")
	}
	if serperKey == "" {
		log.Fatal("SERPER_API_KEY is not set")
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	doMatches := !*perfsOnly
	doPerfs := !*matchesOnly

	if doMatches {
		if err := backfillMatches(ctx, pool, groqKey, serperKey, *dryRun); err != nil {
			log.Fatalf("failed to backfill matches: %v", err)
		}
	}

	if doPerfs {
		if err := backfillPerformances(ctx, pool, groqKey, serperKey, *dryRun); err != nil {
			log.Fatalf("failed to backfill performances: %v", err)
		}
	}

	log.Println("✅ Backfill complete!")
}

// ── Match backfill ──────────────────────────────────────────────────────────

func backfillMatches(ctx context.Context, pool *pgxpool.Pool, groqKey, serperKey string, dryRun bool) error {
	rows, err := pool.Query(ctx, `
		SELECT
			m.id::text,
			m.title,
			ht.name AS home_team_name,
			at.name AS away_team_name,
			m.home_score,
			m.away_score,
			m.home_penalty_score,
			m.away_penalty_score,
			m.round,
			m.venue,
			m.utc_datetime
		FROM matches m
		JOIN teams ht ON m.home_team_id = ht.id
		JOIN teams at ON m.away_team_id = at.id
		ORDER BY m.utc_datetime ASC
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	var matches []matchRow
	for rows.Next() {
		var m matchRow
		if err := rows.Scan(
			&m.ID, &m.Title, &m.HomeTeamName, &m.AwayTeamName,
			&m.HomeScore, &m.AwayScore, &m.HomePenalty, &m.AwayPenalty,
			&m.Round, &m.Venue, &m.UtcDatetime,
		); err != nil {
			return err
		}
		matches = append(matches, m)
	}

	log.Printf("📋 Found %d matches to process", len(matches))

	for i := 0; i < len(matches); i += 1 {
		end := i + 1
		if end > len(matches) {
			end = len(matches)
		}
		batch := matches[i:end]
		batchNum := (i / 1) + 1
		totalBatches := len(matches)

		log.Printf("🏟️  Processing match %d/%d (1 item)...", batchNum, totalBatches)

		// Get search context for the match
		contexts := getMatchContexts(serperKey, batch)

		prompt := buildBatchMatchPrompt(batch, contexts)
		results, err := callGroq(groqKey, prompt)
		if err != nil {
			log.Printf("⚠️  Match %d failed: %v", batchNum, err)
			time.Sleep(5 * time.Second)
			continue
		}

		resultMap := make(map[string]generatedItem)
		for _, r := range results {
			resultMap[r.ID] = r
		}

		for _, m := range batch {
			r, ok := resultMap[m.ID]
			if !ok {
				log.Printf("   ⚠️  No result for match: %s", m.Title)
				continue
			}

			r.Tagline = enforceWordLimit(r.Tagline, 8, 12)
			r.Description = enforceWordLimit(r.Description, 30, 40)

			if dryRun {
				log.Printf("   ✏️  %s", m.Title)
				log.Printf("      Tagline:     %s", r.Tagline)
				log.Printf("      Description: %s", r.Description)
			} else {
				_, err := pool.Exec(ctx,
					`UPDATE matches SET tagline = $2, description = $3, updated_at = now() WHERE id = $1::uuid`,
					m.ID, r.Tagline, r.Description,
				)
				if err != nil {
					log.Printf("   ❌ DB update failed for %s: %v", m.Title, err)
				} else {
					log.Printf("   ✅ %s → %s", m.Title, r.Tagline)
				}
			}
		}
	}

	log.Println("📋 Match backfill done!")
	return nil
}

// ── Performance backfill ────────────────────────────────────────────────────

func backfillPerformances(ctx context.Context, pool *pgxpool.Pool, groqKey, serperKey string, dryRun bool) error {
	rows, err := pool.Query(ctx, `
		SELECT
			p.id::text,
			pl.name AS player_name,
			t.name AS team_name,
			m.title AS match_title,
			ht.name AS home_team_name,
			at.name AS away_team_name,
			m.home_score,
			m.away_score,
			m.round,
			p.average_rating::float8,
			COALESCE((SELECT sv.value::int FROM performance_stats sv JOIN stat_types stt ON stt.id = sv.stat_type_id WHERE sv.performance_id = p.id AND stt.name = 'goals' LIMIT 1), 0) AS goals,
			COALESCE((SELECT sv.value::int FROM performance_stats sv JOIN stat_types stt ON stt.id = sv.stat_type_id WHERE sv.performance_id = p.id AND stt.name = 'assists' LIMIT 1), 0) AS assists,
			p.minutes_played,
			p.is_starter,
			p.captain,
			p.jersey_number
		FROM performances p
		JOIN players pl ON p.player_id = pl.id
		JOIN player_teams pt ON p.player_team_id = pt.id
		JOIN teams t ON pt.team_id = t.id
		JOIN matches m ON p.match_id = m.id
		JOIN teams ht ON m.home_team_id = ht.id
		JOIN teams at ON m.away_team_id = at.id
		ORDER BY m.utc_datetime ASC, p.average_rating DESC NULLS LAST
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	var perfs []perfRow
	for rows.Next() {
		var p perfRow
		if err := rows.Scan(
			&p.ID, &p.PlayerName, &p.TeamName, &p.MatchTitle,
			&p.HomeTeamName, &p.AwayTeamName, &p.HomeScore, &p.AwayScore,
			&p.Round, &p.Rating, &p.Goals, &p.Assists,
			&p.MinutesPlayed, &p.IsStarter, &p.Captain, &p.JerseyNumber,
		); err != nil {
			return err
		}
		perfs = append(perfs, p)
	}

	log.Printf("⚽ Found %d performances to process", len(perfs))

	for i := 0; i < len(perfs); i += 1 {
		end := i + 1
		if end > len(perfs) {
			end = len(perfs)
		}
		batch := perfs[i:end]
		batchNum := (i / 1) + 1
		totalBatches := len(perfs)

		log.Printf("⚽ Processing performance %d/%d (1 item)...", batchNum, totalBatches)

		// Get search context
		contexts := getPerfContexts(serperKey, batch)

		prompt := buildBatchPerfPrompt(batch, contexts)
		results, err := callGroq(groqKey, prompt)
		if err != nil {
			log.Printf("⚠️  Performance %d failed: %v", batchNum, err)
			time.Sleep(5 * time.Second)
			continue
		}

		resultMap := make(map[string]generatedItem)
		for _, r := range results {
			resultMap[r.ID] = r
		}

		for _, p := range batch {
			r, ok := resultMap[p.ID]
			if !ok {
				log.Printf("   ⚠️  No result for: %s (%s)", p.PlayerName, p.MatchTitle)
				continue
			}

			r.Tagline = enforceWordLimit(r.Tagline, 8, 12)
			r.Description = enforceWordLimit(r.Description, 30, 40)

			if dryRun {
				log.Printf("   ✏️  %s (%s)", p.PlayerName, p.MatchTitle)
				log.Printf("      Tagline:     %s", r.Tagline)
				log.Printf("      Description: %s", r.Description)
			} else {
				_, err := pool.Exec(ctx,
					`UPDATE performances SET tagline = $2, description = $3, updated_at = now() WHERE id = $1::uuid`,
					p.ID, r.Tagline, r.Description,
				)
				if err != nil {
					log.Printf("   ❌ DB update failed for %s: %v", p.PlayerName, err)
				} else {
					log.Printf("   ✅ %s → %s", p.PlayerName, r.Tagline)
				}
			}
		}
	}

	log.Println("⚽ Performance backfill done!")
	return nil
}

// ── Serper API caller ───────────────────────────────────────────────────────

func getMatchContexts(apiKey string, matches []matchRow) []string {
	var wg sync.WaitGroup
	contexts := make([]string, len(matches))

	for i, m := range matches {
		wg.Add(1)
		go func(idx int, match matchRow) {
			defer wg.Done()
			q := fmt.Sprintf("FIFA World Cup 2026 %s vs %s match summary goals highlights key moments", match.HomeTeamName, match.AwayTeamName)
			res, err := callSerper(apiKey, q)
			if err == nil && res != "" {
				contexts[idx] = res
			} else {
				contexts[idx] = "No search results available."
			}
		}(i, m)
	}
	wg.Wait()
	return contexts
}

func getPerfContexts(apiKey string, perfs []perfRow) []string {
	var wg sync.WaitGroup
	contexts := make([]string, len(perfs))

	for i, p := range perfs {
		wg.Add(1)
		go func(idx int, perf perfRow) {
			defer wg.Done()
			q := fmt.Sprintf("FIFA World Cup 2026 %s performance stats %s vs %s", perf.PlayerName, perf.HomeTeamName, perf.AwayTeamName)
			res, err := callSerper(apiKey, q)
			if err == nil && res != "" {
				contexts[idx] = res
			} else {
				contexts[idx] = "No specific player search results available."
			}
		}(i, p)
	}
	wg.Wait()
	return contexts
}

func callSerper(apiKey, query string) (string, error) {
	url := "https://google.serper.dev/search"
	reqBody := serperRequest{Q: query}
	
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("X-API-KEY", apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("serper HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var serperResp serperResponse
	if err := json.Unmarshal(body, &serperResp); err != nil {
		return "", err
	}

	var sb strings.Builder
	for i, item := range serperResp.Organic {
		if i >= 3 {
			break // Keep context small, top 3 results
		}
		sb.WriteString(fmt.Sprintf("- %s: %s\n", item.Title, item.Snippet))
	}

	return sb.String(), nil
}

// ── Batch prompt builders ───────────────────────────────────────────────────

func buildBatchMatchPrompt(matches []matchRow, contexts []string) string {
	var sb strings.Builder

	sb.WriteString(`You are a world-class sports journalist writing for PerformX, a premium football analytics platform.

TASK: Generate a tagline and description for the following FIFA World Cup 2026 match.

You are provided with real Google Search results for the match. Incorporate REAL facts from the search snippets into your writing (e.g. goal scorers, key moments).

CRITICAL REQUIREMENT FOR WORD COUNT:
- Tagline: MUST be 8 to 12 words. Punchy, dramatic, like a newspaper headline. Past tense.
- Description: MUST be EXACTLY between 30 and 40 words. Do NOT make it shorter than 30 words or longer than 40 words. Count the words and ensure it meets this condition!
- If the description is too short, expand it by including details about the tournament stage, group standings, stadium atmosphere, the next opponent, or the tactical implications.

STRICT RULES:
1. Do NOT use generic phrases like "a thrilling encounter" or "edge past". Focus on uniqueness.
2. The item MUST include the "id" field EXACTLY as provided — copy-paste it.

EXAMPLES OF GREAT TAGLINES:
- "Mbappé's brace sinks brave Morocco in Round of 16 classic"
- "Penalty drama sends hosts through after goalless stalemate"
- "Five-goal thriller sees Argentina survive Swiss fightback"
- "Clinical Germany dismantle Japan's World Cup dream"
- "Hosts USA roar into quarter-finals with dominant display"

`)

	sb.WriteString("MATCH TO PROCESS:\n\n")
	for i, m := range matches {
		round := "Group Stage"
		if m.Round != nil && *m.Round != "" {
			round = *m.Round
		}
		venue := ""
		if m.Venue != nil {
			venue = *m.Venue
		}
		penalty := ""
		if m.HomePenalty != nil && m.AwayPenalty != nil {
			penalty = fmt.Sprintf(" | Penalties: %d-%d", *m.HomePenalty, *m.AwayPenalty)
		}

		sb.WriteString(fmt.Sprintf("MATCH %d:\n", i+1))
		sb.WriteString(fmt.Sprintf("  id: \"%s\"\n", m.ID))
		sb.WriteString(fmt.Sprintf("  %s vs %s | Score: %d-%d%s\n", m.HomeTeamName, m.AwayTeamName, m.HomeScore, m.AwayScore, penalty))
		sb.WriteString(fmt.Sprintf("  Round: %s | Venue: %s | Date: %s\n", round, venue, m.UtcDatetime.Format("January 2, 2006")))
		sb.WriteString(fmt.Sprintf("  WEB SEARCH CONTEXT:\n%s\n\n", contexts[i]))
	}

	sb.WriteString(fmt.Sprintf(`Return ONLY a valid JSON array with exactly %d objects. No markdown fences, no explanation, no extra text. Just the raw JSON array:
[{"id": "exact-id-from-above", "tagline": "your 8-12 word tagline", "description": "your 30-40 word description"}, ...]`, len(matches)))

	return sb.String()
}

func buildBatchPerfPrompt(perfs []perfRow, contexts []string) string {
	var sb strings.Builder

	sb.WriteString(`You are a world-class sports journalist writing for PerformX, a premium football analytics platform.

TASK: Generate a tagline and description for the following individual player performance in FIFA World Cup 2026.

You are provided with real Google Search results for the player in that match. Incorporate REAL facts from the search snippets into your writing.

CRITICAL REQUIREMENT FOR WORD COUNT:
- Tagline: MUST be 8 to 12 words. Focus on what THIS PLAYER did specifically. Past tense.
- Description: MUST be EXACTLY between 30 and 40 words. Do NOT make it shorter than 30 words or longer than 40 words. Count the words and ensure it meets this condition!
- If the description is too short, expand it by describing the player's role, their impact on the team, key moments, passing details, or significance to their tournament run.

STRICT RULES:
1. Do NOT use generic phrases like "a masterful display" or "showed great quality". Be SPECIFIC to each player based on their stats and context.
2. If a player scored goals, MENTION it in the tagline. If they assisted, mention it.
3. The item MUST include the "id" field EXACTLY as provided — copy-paste it.

EXAMPLES OF GREAT PLAYER TAGLINES:
- "Messi's vision unlocks defense with inch-perfect assist"
- "Bellingham's clinical brace seals quarter-final berth"
- "Sommer stands tall with seven saves in valiant defeat"
- "Ødegaard pulls the strings in commanding midfield display"

`)

	sb.WriteString("PERFORMANCE TO PROCESS:\n\n")
	for i, p := range perfs {
		round := "Group Stage"
		if p.Round != nil && *p.Round != "" {
			round = *p.Round
		}
		rating := 0.0
		if p.Rating != nil {
			rating = *p.Rating
		}
		role := "Substitute"
		if p.IsStarter {
			role = "Starter"
		}
		if p.Captain {
			role += ", Captain"
		}
		statLine := "No goals or assists"
		if p.Goals > 0 && p.Assists > 0 {
			statLine = fmt.Sprintf("%d goal(s), %d assist(s)", p.Goals, p.Assists)
		} else if p.Goals > 0 {
			statLine = fmt.Sprintf("%d goal(s)", p.Goals)
		} else if p.Assists > 0 {
			statLine = fmt.Sprintf("%d assist(s)", p.Assists)
		}

		sb.WriteString(fmt.Sprintf("PERFORMANCE %d:\n", i+1))
		sb.WriteString(fmt.Sprintf("  id: \"%s\"\n", p.ID))
		sb.WriteString(fmt.Sprintf("  Player: %s | Team: %s\n", p.PlayerName, p.TeamName))
		sb.WriteString(fmt.Sprintf("  Match: %s vs %s (%d-%d) | Round: %s\n", p.HomeTeamName, p.AwayTeamName, p.HomeScore, p.AwayScore, round))
		sb.WriteString(fmt.Sprintf("  Role: %s | Minutes: %d | Stats: %s | Rating: %.1f/10\n", role, p.MinutesPlayed, statLine, rating))
		sb.WriteString(fmt.Sprintf("  WEB SEARCH CONTEXT:\n%s\n\n", contexts[i]))
	}

	sb.WriteString(fmt.Sprintf(`Return ONLY a valid JSON array with exactly %d objects. No markdown fences, no explanation, no extra text. Just the raw JSON array:
[{"id": "exact-id-from-above", "tagline": "your 8-12 word tagline", "description": "your 30-40 word description"}, ...]`, len(perfs)))

	return sb.String()
}

// ── Groq API caller ───────────────────────────────────────────────────────

func callGroq(apiKey string, prompt string) ([]generatedItem, error) {
	url := "https://api.groq.com/openai/v1/chat/completions"

	reqBody := groqRequest{
		Model: "llama-3.3-70b-versatile",
		Messages: []groqMessage{
			{Role: "user", Content: prompt},
		},
		Temperature: 0.7,
		MaxTokens:   2000,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	var lastErr error
	for attempt := 0; attempt < 3; attempt++ {
		if attempt > 0 {
			waitSecs := 10 * (attempt + 1)
			log.Printf("      🔄 Retry %d/3, waiting %ds...", attempt+1, waitSecs)
			time.Sleep(time.Duration(waitSecs) * time.Second)
		}

		req, err := http.NewRequest("POST", url, bytes.NewReader(jsonBody))
		if err != nil {
			lastErr = fmt.Errorf("create request: %w", err)
			continue
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+apiKey)

		client := &http.Client{Timeout: 60 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("http request: %w", err)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			lastErr = fmt.Errorf("read body: %w", err)
			continue
		}

		if resp.StatusCode == 429 {
			log.Printf("      ⏳ Rate limited (429), waiting 30s...")
			time.Sleep(30 * time.Second)
			lastErr = fmt.Errorf("rate limited (429)")
			continue
		}

		if resp.StatusCode != 200 {
			lastErr = fmt.Errorf("HTTP %d: %s", resp.StatusCode, truncate(string(body), 400))
			continue
		}

		var groqResp groqResponse
		if err := json.Unmarshal(body, &groqResp); err != nil {
			lastErr = fmt.Errorf("unmarshal response: %w (body: %s)", err, truncate(string(body), 300))
			continue
		}
		if groqResp.Error != nil {
			lastErr = fmt.Errorf("API error: %s", groqResp.Error.Message)
			continue
		}

		fullText := ""
		if len(groqResp.Choices) > 0 {
			fullText = groqResp.Choices[0].Message.Content
		}

		if fullText == "" {
			lastErr = fmt.Errorf("empty response from Groq (body: %s)", truncate(string(body), 400))
			continue
		}

		results, err := extractJSONArray(fullText)
		if err != nil {
			lastErr = fmt.Errorf("extract JSON: %w (text: %s)", err, truncate(fullText, 400))
			continue
		}

		if len(results) == 0 {
			lastErr = fmt.Errorf("Groq returned empty array")
			continue
		}

		return results, nil
	}

	return nil, fmt.Errorf("all retries failed: %w", lastErr)
}

func extractJSONArray(text string) ([]generatedItem, error) {
	cleaned := strings.TrimSpace(text)
	cleaned = strings.TrimPrefix(cleaned, "```json")
	cleaned = strings.TrimPrefix(cleaned, "```JSON")
	cleaned = strings.TrimPrefix(cleaned, "```")
	cleaned = strings.TrimSuffix(cleaned, "```")
	cleaned = strings.TrimSpace(cleaned)

	var results []generatedItem
	if err := json.Unmarshal([]byte(cleaned), &results); err == nil {
		return results, nil
	}

	startIdx := strings.Index(cleaned, "[")
	endIdx := strings.LastIndex(cleaned, "]")
	if startIdx >= 0 && endIdx > startIdx {
		jsonStr := cleaned[startIdx : endIdx+1]
		if parseErr := json.Unmarshal([]byte(jsonStr), &results); parseErr == nil {
			return results, nil
		} else {
			return nil, fmt.Errorf("found array brackets but couldn't parse: %w", parseErr)
		}
	}

	return nil, fmt.Errorf("no JSON array found in response")
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

func enforceWordLimit(text string, min, max int) string {
	words := strings.Fields(text)
	if len(words) > max {
		truncated := strings.Join(words[:max], " ")
		if !strings.HasSuffix(truncated, ".") && !strings.HasSuffix(truncated, "!") && !strings.HasSuffix(truncated, "?") {
			truncated = strings.TrimRight(truncated, ",;:-") + "."
		}
		return truncated
	}
	return text
}

