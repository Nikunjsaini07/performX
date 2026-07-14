package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL must be set")
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer pool.Close()

	fmt.Println("Wiping out old matches, teams, and players for World Cup 2026 pivot...")

	// Truncate tables with CASCADE
	queries := []string{
		"TRUNCATE TABLE performance_review_comments CASCADE;",
		"TRUNCATE TABLE performance_ratings CASCADE;",
		"TRUNCATE TABLE performance_reviews CASCADE;",
		"TRUNCATE TABLE performances CASCADE;",
		"TRUNCATE TABLE match_review_comments CASCADE;",
		"TRUNCATE TABLE match_ratings CASCADE;",
		"TRUNCATE TABLE match_reviews CASCADE;",
		"TRUNCATE TABLE matches CASCADE;",
		"TRUNCATE TABLE player_teams CASCADE;",
		"TRUNCATE TABLE players CASCADE;",
		"TRUNCATE TABLE teams CASCADE;",
	}

	for _, q := range queries {
		_, err := pool.Exec(ctx, q)
		if err != nil {
			log.Fatalf("Error executing query %s: %v", q, err)
		}
	}

	fmt.Println("Database successfully wiped! Ready for FIFA 2026 data.")
}
