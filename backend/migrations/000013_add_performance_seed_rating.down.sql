ALTER TABLE performances DROP COLUMN IF EXISTS seed_rating;

ALTER TABLE performance_ratings ALTER COLUMN rating TYPE NUMERIC(2,1);
ALTER TABLE match_ratings ALTER COLUMN rating TYPE NUMERIC(2,1);
