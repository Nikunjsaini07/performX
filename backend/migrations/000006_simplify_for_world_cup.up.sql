-- 1. Drop List-related tables
DROP TABLE IF EXISTS list_items CASCADE;
DROP TABLE IF EXISTS list_likes CASCADE;
DROP TABLE IF EXISTS lists CASCADE;

-- 2. Drop Follows
DROP TABLE IF EXISTS user_follows CASCADE;

-- 3. Modify Matches to not require season_id
DROP INDEX IF EXISTS matches_season_datetime_idx;
ALTER TABLE matches DROP CONSTRAINT IF EXISTS matches_season_id_fkey;
ALTER TABLE matches DROP COLUMN season_id;

-- Create an index just on utc_datetime now
CREATE INDEX matches_utc_datetime_idx ON matches(utc_datetime);

-- 4. Modify Teams to not require sport_id
DROP INDEX IF EXISTS teams_sport_idx;
ALTER TABLE teams DROP CONSTRAINT IF EXISTS teams_sport_id_fkey;
ALTER TABLE teams DROP COLUMN sport_id;

-- 5. Modify Players to not require sport_id
DROP INDEX IF EXISTS players_sport_id_idx;
ALTER TABLE players DROP CONSTRAINT IF EXISTS players_sport_id_fkey;
ALTER TABLE players DROP COLUMN sport_id;

-- 6. Modify Stat Types to not require sport_id
ALTER TABLE stat_types DROP CONSTRAINT IF EXISTS stat_types_sport_name_unique;
ALTER TABLE stat_types DROP CONSTRAINT IF EXISTS stat_types_sport_short_name_unique;
DROP INDEX IF EXISTS stat_types_sport_category_display_order_idx;

ALTER TABLE stat_types DROP CONSTRAINT IF EXISTS stat_types_sport_id_fkey;
ALTER TABLE stat_types DROP COLUMN sport_id;

ALTER TABLE stat_types ADD CONSTRAINT stat_types_name_unique UNIQUE (name);
ALTER TABLE stat_types ADD CONSTRAINT stat_types_short_name_unique UNIQUE (short_name);
CREATE INDEX stat_types_category_display_order_idx ON stat_types (category, display_order);

-- 7. Drop Hierarchy tables
DROP TABLE IF EXISTS seasons CASCADE;
DROP TABLE IF EXISTS tournaments CASCADE;
DROP TABLE IF EXISTS sports CASCADE;
