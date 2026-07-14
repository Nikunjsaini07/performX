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
	_ = godotenv.Load(".env")
	dbURL := os.Getenv("DATABASE_URL")
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	sqlBytes, err := os.ReadFile("migrations/000006_simplify_for_world_cup.up.sql")
	if err != nil {
		log.Fatal(err)
	}

	_, err = pool.Exec(ctx, string(sqlBytes))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully ran migration #6!")
}
