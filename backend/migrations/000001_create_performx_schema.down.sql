-- ==========================================
-- 1. Drop Social Interaction (Likes & Comments)
-- ==========================================
DROP TABLE IF EXISTS match_review_comment_likes CASCADE;
DROP TABLE IF EXISTS performance_review_comments CASCADE;
DROP TABLE IF EXISTS match_review_comments CASCADE;
DROP TABLE IF EXISTS performance_review_likes CASCADE;
DROP TABLE IF EXISTS match_review_likes CASCADE;
DROP TABLE IF EXISTS performance_reviews CASCADE;
DROP TABLE IF EXISTS match_reviews CASCADE;

-- ==========================================
-- 2. Drop User Content (Lists, Ratings, Follows)
-- ==========================================
DROP TABLE IF EXISTS user_follows CASCADE;
DROP TABLE IF EXISTS list_likes CASCADE;
DROP TABLE IF EXISTS list_items CASCADE;
DROP TABLE IF EXISTS lists CASCADE;
DROP TABLE IF EXISTS performance_ratings CASCADE;
DROP TABLE IF EXISTS match_ratings CASCADE;

-- ==========================================
-- 3. Drop Player Performance & Stats
-- ==========================================
DROP TABLE IF EXISTS performance_stats CASCADE;
DROP TABLE IF EXISTS match_stats CASCADE;
DROP TABLE IF EXISTS stat_types CASCADE;
DROP TABLE IF EXISTS performances CASCADE;

-- ==========================================
-- 4. Drop Match Context
-- ==========================================
DROP TABLE IF EXISTS matches CASCADE;

-- ==========================================
-- 5. Drop Domain Context (Teams, Seasons, Tournaments, Countries, Sports, Users)
-- ==========================================
DROP TABLE IF EXISTS player_teams CASCADE;
DROP TABLE IF EXISTS players CASCADE;
DROP TABLE IF EXISTS teams CASCADE;
DROP TABLE IF EXISTS seasons CASCADE;
DROP TABLE IF EXISTS tournaments CASCADE;
DROP TABLE IF EXISTS countries CASCADE;
DROP TABLE IF EXISTS sports CASCADE;
DROP TABLE IF EXISTS users CASCADE;

-- ==========================================
-- 6. Drop Enums and Extensions
-- ==========================================
DROP TYPE IF EXISTS team_type CASCADE;
DROP TYPE IF EXISTS tournament_type CASCADE;
DROP TYPE IF EXISTS review_status CASCADE;
DROP TYPE IF EXISTS user_role CASCADE;
