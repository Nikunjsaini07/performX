-- Preserve the original provider (FotMob) rating as a permanent seed so the
-- community average can always be recomputed as:
--   average_rating = (seed_rating + SUM(user ratings)) / (1 + COUNT(user ratings))
-- keeping the stored column (homepage) and the live detail-page query in sync.

ALTER TABLE performances ADD COLUMN IF NOT EXISTS seed_rating NUMERIC(3,1);

-- Widen the rating columns: numeric(2,1) could only store up to 9.9 even though
-- the CHECK constraint permits a perfect 10.0.
ALTER TABLE performance_ratings ALTER COLUMN rating TYPE NUMERIC(3,1);
ALTER TABLE match_ratings ALTER COLUMN rating TYPE NUMERIC(3,1);

-- Backfill seed_rating.
-- Performances with no community votes yet (total_votes = 1) still hold the
-- provider rating in average_rating, so copy it directly.
UPDATE performances
SET seed_rating = average_rating
WHERE seed_rating IS NULL AND (total_votes IS NULL OR total_votes <= 1);

-- The launch-seeded performances (total_votes > 1) had average_rating blended,
-- so restore their exact original provider seeds by slug.
UPDATE performances SET seed_rating = 9.7 WHERE slug = 'lionel-messi-c675-vs-algeria-763d920d';
UPDATE performances SET seed_rating = 9.6 WHERE slug = 'jonathan-david-5e39-vs-qatar-82ba58b8';
UPDATE performances SET seed_rating = 9.6 WHERE slug = 'ousmane-demb-l-be95-vs-norway-90b387d7';
UPDATE performances SET seed_rating = 9.4 WHERE slug = 'charles-de-ketelaere-c8aa-vs-usa-fe4e0dea';
UPDATE performances SET seed_rating = 9.4 WHERE slug = 'jude-bellingham-vs-norway-so5srw';
UPDATE performances SET seed_rating = 9.4 WHERE slug = 'kylian-mbapp-cb77-vs-sweden-118035bc';
UPDATE performances SET seed_rating = 9.3 WHERE slug = 'mikel-oyarzabal-a7c3-vs-saudi-arabia-279609dc';
UPDATE performances SET seed_rating = 9.3 WHERE slug = 'ayase-ueda-260e-vs-tunisia-964a9fa6';
UPDATE performances SET seed_rating = 9.3 WHERE slug = 'vin-cius-j-nior-bce7-vs-scotland-78b8b109';
UPDATE performances SET seed_rating = 9.3 WHERE slug = 'eloy-room-a95d-vs-ecuador-355b37de';
UPDATE performances SET seed_rating = 9.3 WHERE slug = 'pape-gueye-c629-vs-iraq-c4ccfafb';
UPDATE performances SET seed_rating = 9.3 WHERE slug = 'cody-gakpo-eaf0-vs-sweden-04c1fe2e';

-- Any remaining nulls (safety net): fall back to current average_rating.
UPDATE performances SET seed_rating = COALESCE(average_rating, 0.0) WHERE seed_rating IS NULL;
