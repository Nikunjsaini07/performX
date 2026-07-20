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
	"github.com/Nikunjsaini07/performx/backend/internal/middleware"
	
)

func loadEnv() (string, string, string) {
	var dbURL, jwtSecret, port string
	port = "8080" // default

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
			
			
			os.Setenv(key, val)
		}
	}

	dbURL = os.Getenv("DATABASE_URL")
	jwtSecret = os.Getenv("JWT_SECRET")
	
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	return dbURL, jwtSecret, port
}

func main() {
	dbURL, jwtSecret, port := loadEnv()
	if dbURL == "" {
		log.Fatal("DATABASE_URL is required in .env or OS environment variables")
	}
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is required in .env or OS environment variables")
	}

	ctx := context.Background()
	
	
	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		log.Fatalf("Unable to parse database URL: %v", err)
	}

	
	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnIdleTime = 15 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatalf("Unable to connect to database pool: %v", err)
	}
	defer pool.Close()

	
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Database connection ping failed: %v", err)
	}
	fmt.Println("🚀 Connected to Neon PostgreSQL database pool successfully!")

	queries := db.New(pool)

	
	mux := http.NewServeMux()
	api.RegisterRoutes(mux, queries, []byte(jwtSecret))

	

	serverAddr := fmt.Sprintf(":%s", port)
	server := &http.Server{
		Addr:         serverAddr,
		Handler:      middleware.EnableCORS(mux),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	fmt.Printf("🔥 PerformX Auth API Server is running on port %s...\n", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server ListenAndServe failed: %v", err)
	}
}
