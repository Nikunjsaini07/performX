package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load(".env")
	dbURL := os.Getenv("DATABASE_URL")
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	_, err = pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version TEXT PRIMARY KEY
		);
	`)
	if err != nil {
		log.Fatalf("Failed to create schema_migrations table: %v", err)
	}

	files, err := os.ReadDir("migrations")
	if err != nil {
		log.Fatalf("Failed to read migrations directory: %v", err)
	}

	var upMigrations []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".up.sql") {
			upMigrations = append(upMigrations, file.Name())
		}
	}
	sort.Strings(upMigrations)

	var count int
	err = pool.QueryRow(ctx, "SELECT COUNT(*) FROM schema_migrations").Scan(&count)
	if err == nil && count == 0 {
		var usersExists bool
		err = pool.QueryRow(ctx, "SELECT EXISTS (SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'users')").Scan(&usersExists)
		
		if err == nil && usersExists {
			fmt.Println("⚠️ Pre-existing production database detected!")
			fmt.Println("Backfilling migration history instead of executing...")
			for _, filename := range upMigrations {
				_, err = pool.Exec(ctx, "INSERT INTO schema_migrations (version) VALUES ($1)", filename)
				if err != nil {
					log.Fatalf("Failed to backfill %s: %v", filename, err)
				}
			}
			fmt.Println("✅ Backfill complete. You are now safe to add new migrations in the future.")
			return
		}
	}

	for _, filename := range upMigrations {
		var exists bool
		err := pool.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version = $1)", filename).Scan(&exists)
		if err != nil {
			log.Fatalf("Failed to check migration status: %v", err)
		}

		if exists {
			fmt.Printf("Skipping %s (already applied)\n", filename)
			continue
		}

		fmt.Printf("Applying %s...\n", filename)
		
		path := filepath.Join("migrations", filename)
		sqlBytes, err := os.ReadFile(path)
		if err != nil {
			log.Fatalf("Failed to read %s: %v", filename, err)
		}

		_, err = pool.Exec(ctx, string(sqlBytes))
		if err != nil {
			log.Fatalf("Failed to execute %s: %v", filename, err)
		}

		_, err = pool.Exec(ctx, "INSERT INTO schema_migrations (version) VALUES ($1)", filename)
		if err != nil {
			log.Fatalf("Failed to record %s in schema_migrations: %v", filename, err)
		}
	}

	fmt.Println("🎉 All migrations applied successfully!")
}
