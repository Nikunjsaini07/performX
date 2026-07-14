package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}
	dbUrl := os.Getenv("DATABASE_URL")
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, dbUrl)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer conn.Close(ctx)

	// 1. Insert 4 Test Users
	fmt.Println("Inserting test users...")
	usersQuery := `
		INSERT INTO users (username, display_name, email, bio, avatar_url) VALUES 
		('alex_football', 'Alex (Football Nerd)', 'alex@example.com', 'I watch every single game.', 'https://api.dicebear.com/7.x/avataaars/svg?seed=Alex&backgroundColor=b6e3f4'),
		('sarah_stats', 'Sarah Stats', 'sarah@example.com', 'Data analyst by day, football fan by night.', 'https://api.dicebear.com/7.x/avataaars/svg?seed=Sarah&backgroundColor=ffdfbf'),
		('messi_goat', 'Leo Fan 10', 'messi_goat@example.com', 'Visca el Barca.', 'https://api.dicebear.com/7.x/avataaars/svg?seed=Leo&backgroundColor=c0aede'),
		('tactic_master', 'Tactics Master', 'tactic@example.com', 'It is all about the low block.', 'https://api.dicebear.com/7.x/avataaars/svg?seed=Tactic&backgroundColor=d1d4f9')
		ON CONFLICT (username) DO NOTHING;
	`
	if _, err := conn.Exec(ctx, usersQuery); err != nil {
		log.Fatalf("Error inserting users: %v", err)
	}

	// 2. Update Player Photo URLs
	fmt.Println("Updating player portraits...")
	portraits := map[string]string{
		"lionel-messi":       "https://images.unsplash.com/photo-1698263051016-04d166dfbc08?q=80&w=2787&auto=format&fit=crop", // placeholder portrait
		"cristiano-ronaldo":  "https://images.unsplash.com/photo-1658428801550-98cc76eb73b0?q=80&w=2835&auto=format&fit=crop", // placeholder portrait
		"kylian-mbappe":      "https://images.unsplash.com/photo-1678122998632-44675e01239c?q=80&w=2787&auto=format&fit=crop", // placeholder portrait
		"jude-bellingham":    "https://images.unsplash.com/photo-1695629165971-d6fcda28a0ed?q=80&w=2864&auto=format&fit=crop", // placeholder portrait
		"kevin-de-bruyne":    "https://images.unsplash.com/photo-1543326727-cf6c39e8f84c?q=80&w=2940&auto=format&fit=crop",
		"erling-haaland":     "https://images.unsplash.com/photo-1579952363873-27f3bade9f55?q=80&w=2835&auto=format&fit=crop",
		"vinicius-junior":    "https://images.unsplash.com/photo-1614632537423-1e6c2e7e0aab?q=80&w=2940&auto=format&fit=crop",
		"mohamed-salah":      "https://images.unsplash.com/photo-1517466787929-bc90951d0974?q=80&w=2786&auto=format&fit=crop",
	}

	for slug, url := range portraits {
		_, err := conn.Exec(ctx, "UPDATE players SET photo_url = $1 WHERE slug = $2", url, slug)
		if err != nil {
			log.Printf("Failed to update photo for %s: %v\n", slug, err)
		}
	}

	// 3. Insert Test Reviews for Trending Feed
	fmt.Println("Generating mock reviews & ratings...")
	reviewsQuery := `
		DO $$ 
		DECLARE 
			user1_id UUID;
			user2_id UUID;
			user3_id UUID;
			user4_id UUID;
			perf_record RECORD;
			rating1 NUMERIC(2,1);
			rating2 NUMERIC(2,1);
			rating3 NUMERIC(2,1);
		BEGIN
			SELECT id INTO user1_id FROM users WHERE username = 'alex_football' LIMIT 1;
			SELECT id INTO user2_id FROM users WHERE username = 'sarah_stats' LIMIT 1;
			SELECT id INTO user3_id FROM users WHERE username = 'messi_goat' LIMIT 1;
			SELECT id INTO user4_id FROM users WHERE username = 'tactic_master' LIMIT 1;
			
			FOR perf_record IN SELECT id FROM performances LIMIT 20 LOOP
				IF random() > 0.5 THEN
					rating1 := floor(random() * (9 - 7 + 1) + 7);
					INSERT INTO performance_ratings (performance_id, user_id, rating)
					VALUES (perf_record.id, user1_id, rating1)
					ON CONFLICT DO NOTHING;

					INSERT INTO performance_reviews (performance_id, user_id, content)
					VALUES (perf_record.id, user1_id, 'Incredible performance, one for the history books!')
					ON CONFLICT DO NOTHING;
				END IF;
				
				IF random() > 0.6 THEN
					rating2 := floor(random() * (9 - 8 + 1) + 8);
					INSERT INTO performance_ratings (performance_id, user_id, rating)
					VALUES (perf_record.id, user2_id, rating2)
					ON CONFLICT DO NOTHING;

					INSERT INTO performance_reviews (performance_id, user_id, content)
					VALUES (perf_record.id, user2_id, 'Statistically brilliant. Highest xG contribution of the match.')
					ON CONFLICT DO NOTHING;
				END IF;

				IF random() > 0.7 THEN
					rating3 := floor(random() * (8 - 5 + 1) + 5);
					INSERT INTO performance_ratings (performance_id, user_id, rating)
					VALUES (perf_record.id, user3_id, rating3)
					ON CONFLICT DO NOTHING;

					INSERT INTO performance_reviews (performance_id, user_id, content)
					VALUES (perf_record.id, user3_id, 'Honestly expected more, but still a solid display.')
					ON CONFLICT DO NOTHING;
				END IF;
			END LOOP;
		END $$;
	`
	if _, err := conn.Exec(ctx, reviewsQuery); err != nil {
		log.Printf("Error inserting reviews: %v", err)
	}

	fmt.Println("Seed data generated successfully!")
}
