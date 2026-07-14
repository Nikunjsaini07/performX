package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Nikunjsaini07/performx/backend/internal/api"
	"github.com/Nikunjsaini07/performx/backend/internal/db"
	"github.com/Nikunjsaini07/performx/backend/internal/worker"
)

// Helper to parse environment variables from .env file
func loadEnv() (string, string, string, string) {
	var dbURL, jwtSecret, port, adminPrefix string
	port = "8080" // default
	jwtSecret = "performx-default-super-secret-key-change-in-prod"
	adminPrefix = "admin-gate-performx"

	content, err := os.ReadFile(".env")
	if err == nil {
		lines := strings.Split(string(content), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			parts := strings.SplitN(line, "=", 2)
			if len(parts) != 2 {
				continue
			}
			key := strings.TrimSpace(parts[0])
			val := strings.Trim(strings.TrimSpace(parts[1]), `'"`)
			
			// Set as environment variable so os.Getenv() works throughout the app
			os.Setenv(key, val)
			
			switch key {
			case "DATABASE_URL":
				dbURL = val
			case "JWT_SECRET":
				jwtSecret = val
			case "PORT":
				port = val
			case "ADMIN_ROUTE_PREFIX":
				adminPrefix = val
			}
		}
	}

	// Fallback to OS environment
	if dbURL == "" {
		dbURL = os.Getenv("DATABASE_URL")
	}
	if os.Getenv("JWT_SECRET") != "" {
		jwtSecret = os.Getenv("JWT_SECRET")
	}
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	if os.Getenv("ADMIN_ROUTE_PREFIX") != "" {
		adminPrefix = os.Getenv("ADMIN_ROUTE_PREFIX")
	}

	return dbURL, jwtSecret, port, adminPrefix
}

func main() {
	dbURL, jwtSecret, port, adminPrefix := loadEnv()
	if dbURL == "" {
		log.Fatal("DATABASE_URL is required in .env or OS environment variables")
	}

	ctx := context.Background()
	
	// Create database connection pool
	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		log.Fatalf("Unable to parse database URL: %v", err)
	}

	// Tweak pool settings for performance
	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnIdleTime = 15 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatalf("Unable to connect to database pool: %v", err)
	}
	defer pool.Close()

	// Verify database connection is healthy
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Database connection ping failed: %v", err)
	}
	fmt.Println("🚀 Connected to Neon PostgreSQL database pool successfully!")

	queries := db.New(pool)

	// Setup routes using Go 1.22+ enhanced ServeMux router
	mux := http.NewServeMux()
	api.RegisterRoutes(mux, queries, []byte(jwtSecret), adminPrefix)

	// Start the background worker for trending scores
	worker.StartTrendingWorker(pool, queries, 30*time.Minute)

	serverAddr := fmt.Sprintf(":%s", port)
	server := &http.Server{
		Addr:         serverAddr,
		Handler:      api.EnableCORS(mux),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	fmt.Printf("🔥 PerformX Auth API Server is running on port %s...\n", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server ListenAndServe failed: %v", err)
	}
}
