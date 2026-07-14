ALTER TABLE performances RENAME COLUMN provider_rating TO average_rating;
ALTER TABLE performances ADD COLUMN total_votes INT DEFAULT 1;

ALTER TABLE matches ADD COLUMN average_rating DECIMAL(3,1) DEFAULT 0.0;
ALTER TABLE matches ADD COLUMN total_votes INT DEFAULT 1;
