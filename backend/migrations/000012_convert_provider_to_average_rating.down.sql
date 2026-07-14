ALTER TABLE performances RENAME COLUMN average_rating TO provider_rating;
ALTER TABLE performances DROP COLUMN total_votes;

ALTER TABLE matches RENAME COLUMN average_rating TO provider_rating;
ALTER TABLE matches DROP COLUMN total_votes;
