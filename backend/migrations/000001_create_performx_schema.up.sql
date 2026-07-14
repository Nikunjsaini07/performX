CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TYPE user_role AS ENUM ('USER', 'ADMIN');
CREATE TYPE review_status AS ENUM ('PENDING', 'APPROVED', 'FLAGGED', 'REMOVED');
CREATE TYPE tournament_type AS ENUM (
    'LEAGUE',
    'CUP'
);
CREATE TYPE team_type AS ENUM (
    'CLUB',
    'NATIONAL'
);

-- ==========================================
-- 1. Independent Core Tables
-- ==========================================

CREATE TABLE sports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    slug TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT sports_name_unique UNIQUE (name),
    CONSTRAINT sports_slug_unique UNIQUE (slug),
    CONSTRAINT sports_name_not_blank CHECK (btrim(name) <> ''),
    CONSTRAINT sports_slug_not_blank CHECK (btrim(slug) <> '')
);

CREATE TABLE countries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    iso2 CHAR(2) NOT NULL,
    iso3 CHAR(3) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT countries_name_unique UNIQUE (name),
    CONSTRAINT countries_iso2_unique UNIQUE (iso2),
    CONSTRAINT countries_iso3_unique UNIQUE (iso3),
    CONSTRAINT countries_name_not_blank CHECK (btrim(name) <> ''),
    CONSTRAINT countries_iso2_uppercase CHECK (iso2 = upper(iso2)),
    CONSTRAINT countries_iso3_uppercase CHECK (iso3 = upper(iso3))
);

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username TEXT NOT NULL,
    display_name TEXT NOT NULL,
    email TEXT NOT NULL,
    bio TEXT,
    avatar_url TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT users_username_unique UNIQUE (username),
    CONSTRAINT users_email_unique UNIQUE (email),
    CONSTRAINT users_username_not_blank CHECK (btrim(username) <> ''),
    CONSTRAINT users_display_name_not_blank CHECK (btrim(display_name) <> ''),
    CONSTRAINT users_email_not_blank CHECK (btrim(email) <> ''),
    CONSTRAINT users_bio_not_blank CHECK (bio IS NULL OR btrim(bio) <> '')
);

-- ==========================================
-- 2. Domain Hierarchy (Sports Context)
-- ==========================================

CREATE TABLE tournaments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sport_id UUID NOT NULL REFERENCES sports(id) ON DELETE RESTRICT,
    name TEXT NOT NULL,
    short_name TEXT NOT NULL,
    slug TEXT NOT NULL,
    type tournament_type NOT NULL,
    logo_url TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT tournaments_name_unique UNIQUE (name),
    CONSTRAINT tournaments_short_name_unique UNIQUE (short_name),
    CONSTRAINT tournaments_slug_unique UNIQUE (slug),
    CONSTRAINT tournaments_name_not_blank CHECK (btrim(name) <> ''),
    CONSTRAINT tournaments_short_name_not_blank CHECK (btrim(short_name) <> ''),
    CONSTRAINT tournaments_slug_not_blank CHECK (btrim(slug) <> '')
);

CREATE INDEX tournaments_sport_id_idx ON tournaments(sport_id);

CREATE TABLE seasons (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tournament_id UUID NOT NULL REFERENCES tournaments(id) ON DELETE RESTRICT,
    name TEXT NOT NULL,
    start_year INTEGER NOT NULL,
    end_year INTEGER NOT NULL,
    start_date DATE,
    end_date DATE,
    is_current BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT seasons_tournament_name_unique UNIQUE (tournament_id, name),
    CONSTRAINT seasons_tournament_years_unique UNIQUE (tournament_id, start_year, end_year),
    CONSTRAINT seasons_name_not_blank CHECK (btrim(name) <> ''),
    CONSTRAINT seasons_year_range_valid CHECK (end_year >= start_year),
    CONSTRAINT seasons_date_range_valid CHECK (
        start_date IS NULL
        OR end_date IS NULL
        OR end_date >= start_date
    ),
    CONSTRAINT seasons_years_reasonable CHECK (
        start_year BETWEEN 1800 AND 3000
        AND end_year BETWEEN 1800 AND 3000
    )
);

CREATE INDEX seasons_tournament_idx ON seasons(tournament_id);

CREATE TABLE teams (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sport_id UUID NOT NULL REFERENCES sports(id) ON DELETE RESTRICT,
    country_id UUID NOT NULL REFERENCES countries(id) ON DELETE RESTRICT,
    name TEXT NOT NULL,
    short_name TEXT NOT NULL,
    slug TEXT NOT NULL,
    type team_type NOT NULL,
    logo_url TEXT,
    founded_year INTEGER,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT teams_slug_unique UNIQUE (slug),
    CONSTRAINT teams_name_country_unique UNIQUE (name, country_id),
    CONSTRAINT teams_short_name_country_unique UNIQUE (short_name, country_id),
    CONSTRAINT teams_name_not_blank CHECK (btrim(name) <> ''),
    CONSTRAINT teams_short_name_not_blank CHECK (btrim(short_name) <> ''),
    CONSTRAINT teams_slug_not_blank CHECK (btrim(slug) <> ''),
    CONSTRAINT teams_founded_year_reasonable CHECK (
        founded_year IS NULL
        OR founded_year BETWEEN 1800 AND 3000
    )
);

CREATE INDEX teams_country_idx ON teams(country_id);
CREATE INDEX teams_sport_idx ON teams(sport_id);

CREATE TABLE players (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sport_id UUID NOT NULL REFERENCES sports(id) ON DELETE RESTRICT,
    name TEXT NOT NULL,
    slug TEXT NOT NULL,
    full_name TEXT,
    known_as TEXT,
    date_of_birth DATE,
    place_of_birth TEXT,
    country_id UUID REFERENCES countries(id) ON DELETE RESTRICT,
    photo_url TEXT,
    height_cm INTEGER,
    weight_kg INTEGER,
    shirt_name TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT players_slug_unique UNIQUE (slug),
    CONSTRAINT players_name_not_blank CHECK (btrim(name) <> ''),
    CONSTRAINT players_slug_not_blank CHECK (btrim(slug) <> ''),
    CONSTRAINT players_full_name_not_blank CHECK (full_name IS NULL OR btrim(full_name) <> ''),
    CONSTRAINT players_known_as_not_blank CHECK (known_as IS NULL OR btrim(known_as) <> ''),
    CONSTRAINT players_place_of_birth_not_blank CHECK (place_of_birth IS NULL OR btrim(place_of_birth) <> ''),
    CONSTRAINT players_height_cm_reasonable CHECK (
        height_cm IS NULL
        OR height_cm BETWEEN 100 AND 250
    ),
    CONSTRAINT players_weight_kg_reasonable CHECK (
        weight_kg IS NULL
        OR weight_kg BETWEEN 30 AND 150
    ),
    CONSTRAINT players_shirt_name_not_blank CHECK (
        shirt_name IS NULL
        OR btrim(shirt_name) <> ''
    )
);

CREATE INDEX players_country_id_idx ON players(country_id);
CREATE INDEX players_sport_id_idx ON players(sport_id);

CREATE TABLE player_teams (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    team_id UUID NOT NULL REFERENCES teams(id) ON DELETE RESTRICT,
    jersey_number INTEGER,
    start_date DATE NOT NULL,
    end_date DATE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,

    CONSTRAINT player_teams_player_team_start_unique UNIQUE (player_id, team_id, start_date),
    CONSTRAINT player_teams_id_player_unique UNIQUE (id, player_id),
    CONSTRAINT player_teams_jersey_number_valid CHECK (
        jersey_number IS NULL
        OR jersey_number BETWEEN 1 AND 99
    ),
    CONSTRAINT player_teams_date_range_valid CHECK (
        end_date IS NULL
        OR end_date >= start_date
    )
);

CREATE INDEX player_teams_player_id_idx ON player_teams(player_id);
CREATE INDEX player_teams_team_id_idx ON player_teams(team_id);

-- ==========================================
-- 3. Match Context
-- ==========================================

CREATE TABLE matches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    season_id UUID NOT NULL REFERENCES seasons(id) ON DELETE RESTRICT,
    home_team_id UUID NOT NULL REFERENCES teams(id) ON DELETE RESTRICT,
    away_team_id UUID NOT NULL REFERENCES teams(id) ON DELETE RESTRICT,
    title TEXT NOT NULL,
    slug TEXT NOT NULL,
    description TEXT,
    round TEXT,
    utc_datetime TIMESTAMPTZ NOT NULL,
    venue TEXT,
    cover_image_url TEXT,
    home_score INTEGER NOT NULL,
    away_score INTEGER NOT NULL,
    home_penalty_score INTEGER,
    away_penalty_score INTEGER,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT matches_slug_unique UNIQUE (slug),
    CONSTRAINT matches_home_away_different CHECK (home_team_id <> away_team_id),
    CONSTRAINT matches_title_not_blank CHECK (btrim(title) <> ''),
    CONSTRAINT matches_slug_not_blank CHECK (btrim(slug) <> ''),
    CONSTRAINT matches_description_not_blank CHECK (description IS NULL OR btrim(description) <> ''),
    CONSTRAINT matches_round_not_blank CHECK (round IS NULL OR btrim(round) <> ''),
    CONSTRAINT matches_home_score_valid CHECK (home_score >= 0),
    CONSTRAINT matches_away_score_valid CHECK (away_score >= 0),
    CONSTRAINT matches_home_penalty_score_valid CHECK (
        home_penalty_score IS NULL
        OR home_penalty_score >= 0
    ),
    CONSTRAINT matches_away_penalty_score_valid CHECK (
        away_penalty_score IS NULL
        OR away_penalty_score >= 0
    )
);

CREATE INDEX matches_season_datetime_idx ON matches(season_id, utc_datetime);
CREATE INDEX matches_home_team_idx ON matches(home_team_id);
CREATE INDEX matches_away_team_idx ON matches(away_team_id);

CREATE TABLE match_ratings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    match_id UUID NOT NULL REFERENCES matches(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    rating NUMERIC(2,1) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT match_ratings_user_match_unique UNIQUE (match_id, user_id),
    CONSTRAINT match_ratings_rating_valid CHECK (rating BETWEEN 0 AND 10)
);

CREATE INDEX match_ratings_match_idx ON match_ratings(match_id);
CREATE INDEX match_ratings_user_idx ON match_ratings(user_id);

-- ==========================================
-- 4. Player Performance & Stats
-- ==========================================

CREATE TABLE performances (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    match_id UUID NOT NULL REFERENCES matches(id) ON DELETE CASCADE,
    player_id UUID NOT NULL REFERENCES players(id) ON DELETE RESTRICT,
    player_team_id UUID NOT NULL REFERENCES player_teams(id) ON DELETE RESTRICT,
    title TEXT NOT NULL,
    description TEXT,
    cover_image_url TEXT,
    jersey_number INTEGER,
    is_starter BOOLEAN NOT NULL DEFAULT FALSE,
    captain BOOLEAN NOT NULL DEFAULT FALSE,
    minutes_played INTEGER NOT NULL DEFAULT 0,
    provider_rating NUMERIC(3,1),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT performances_match_player_unique UNIQUE (match_id, player_id),
    CONSTRAINT performances_player_team_matches_player_fk FOREIGN KEY (player_team_id, player_id) REFERENCES player_teams(id, player_id) ON DELETE RESTRICT,
    CONSTRAINT performances_title_not_blank CHECK (btrim(title) <> ''),
    CONSTRAINT performances_description_not_blank CHECK (description IS NULL OR btrim(description) <> ''),
    CONSTRAINT performances_jersey_number_valid CHECK (
        jersey_number IS NULL
        OR jersey_number BETWEEN 1 AND 99
    ),
    CONSTRAINT performances_minutes_played_valid CHECK (minutes_played BETWEEN 0 AND 130),
    CONSTRAINT performances_provider_rating_valid CHECK (
        provider_rating IS NULL
        OR provider_rating BETWEEN 0 AND 10
    )
);

CREATE INDEX performances_match_id_idx ON performances(match_id);
CREATE INDEX performances_player_id_idx ON performances(player_id);
CREATE INDEX performances_player_team_id_idx ON performances(player_team_id);

CREATE TABLE stat_types (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sport_id UUID NOT NULL REFERENCES sports(id) ON DELETE RESTRICT,
    name TEXT NOT NULL,
    short_name TEXT NOT NULL,
    unit TEXT,
    category TEXT NOT NULL,
    display_order INTEGER NOT NULL DEFAULT 0,

    CONSTRAINT stat_types_sport_name_unique UNIQUE (sport_id, name),
    CONSTRAINT stat_types_sport_short_name_unique UNIQUE (sport_id, short_name),
    CONSTRAINT stat_types_name_not_blank CHECK (btrim(name) <> ''),
    CONSTRAINT stat_types_short_name_not_blank CHECK (btrim(short_name) <> ''),
    CONSTRAINT stat_types_category_not_blank CHECK (btrim(category) <> '')
);

CREATE INDEX stat_types_sport_category_display_order_idx ON stat_types (sport_id, category, display_order);

CREATE TABLE match_stats (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    match_id UUID NOT NULL REFERENCES matches(id) ON DELETE CASCADE,
    team_id UUID NOT NULL REFERENCES teams(id) ON DELETE RESTRICT,
    stat_type_id UUID NOT NULL REFERENCES stat_types(id) ON DELETE RESTRICT,
    value NUMERIC(12,3) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT match_stats_match_team_stat_unique UNIQUE (match_id, team_id, stat_type_id)
);

CREATE INDEX match_stats_match_team_idx ON match_stats(match_id, team_id);
CREATE INDEX match_stats_stat_type_value_idx ON match_stats(stat_type_id, value DESC);

CREATE TABLE performance_stats (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    performance_id UUID NOT NULL REFERENCES performances(id) ON DELETE CASCADE,
    stat_type_id UUID NOT NULL REFERENCES stat_types(id) ON DELETE RESTRICT,
    value NUMERIC(12,3) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT performance_stats_performance_stat_type_unique UNIQUE (performance_id, stat_type_id)
);

CREATE INDEX performance_stats_stat_type_value_idx ON performance_stats (stat_type_id, value DESC);
CREATE INDEX performance_stats_performance_idx ON performance_stats (performance_id);

-- ==========================================
-- 5. User Content (Lists, Ratings, Follows)
-- ==========================================

CREATE TABLE lists (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    slug TEXT NOT NULL,
    description TEXT,
    cover_image_url TEXT,
    is_public BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT lists_slug_unique UNIQUE (slug),
    CONSTRAINT lists_title_not_blank CHECK (btrim(title) <> ''),
    CONSTRAINT lists_slug_not_blank CHECK (btrim(slug) <> ''),
    CONSTRAINT lists_description_not_blank CHECK (description IS NULL OR btrim(description) <> '')
);

CREATE INDEX lists_user_idx ON lists(user_id);
CREATE INDEX lists_public_idx ON lists(is_public);

CREATE TABLE list_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    list_id UUID NOT NULL REFERENCES lists(id) ON DELETE CASCADE,
    match_id UUID REFERENCES matches(id) ON DELETE CASCADE,
    performance_id UUID REFERENCES performances(id) ON DELETE CASCADE,
    position INTEGER NOT NULL,
    note TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT list_items_one_entity CHECK (
        (match_id IS NOT NULL AND performance_id IS NULL)
        OR
        (match_id IS NULL AND performance_id IS NOT NULL)
    ),
    CONSTRAINT list_items_position_positive CHECK (position > 0)
);

CREATE INDEX list_items_list_idx ON list_items(list_id);
CREATE INDEX list_items_match_idx ON list_items(match_id);
CREATE INDEX list_items_performance_idx ON list_items(performance_id);

CREATE TABLE list_likes (
    list_id UUID NOT NULL REFERENCES lists(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    PRIMARY KEY (list_id, user_id)
);

CREATE INDEX list_likes_user_idx ON list_likes(user_id);

CREATE TABLE performance_ratings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    performance_id UUID NOT NULL REFERENCES performances(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    rating NUMERIC(2,1) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT performance_ratings_user_performance_unique UNIQUE (performance_id, user_id),
    CONSTRAINT performance_ratings_rating_valid CHECK (rating BETWEEN 0 AND 10)
);

CREATE INDEX performance_ratings_performance_idx ON performance_ratings(performance_id);
CREATE INDEX performance_ratings_user_idx ON performance_ratings(user_id);

CREATE TABLE user_follows (
    follower_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    following_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    PRIMARY KEY (follower_id, following_id),
    CONSTRAINT user_follows_no_self_follow CHECK (follower_id <> following_id)
);

CREATE INDEX user_follows_following_idx ON user_follows(following_id);

-- ==========================================
-- 6. Social Interaction (Reviews & Comments)
-- ==========================================

CREATE TABLE match_reviews (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    match_id UUID NOT NULL REFERENCES matches(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title TEXT,
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT match_reviews_user_match_unique UNIQUE (match_id, user_id),
    CONSTRAINT match_reviews_content_not_blank CHECK (btrim(content) <> ''),
    CONSTRAINT match_reviews_title_not_blank CHECK (title IS NULL OR btrim(title) <> '')
);

CREATE INDEX match_reviews_match_idx ON match_reviews(match_id);
CREATE INDEX match_reviews_user_idx ON match_reviews(user_id);

CREATE TABLE performance_reviews (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    performance_id UUID NOT NULL REFERENCES performances(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title TEXT,
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT performance_reviews_user_performance_unique UNIQUE (performance_id, user_id),
    CONSTRAINT performance_reviews_content_not_blank CHECK (btrim(content) <> ''),
    CONSTRAINT performance_reviews_title_not_blank CHECK (title IS NULL OR btrim(title) <> '')
);

CREATE INDEX performance_reviews_performance_idx ON performance_reviews(performance_id);
CREATE INDEX performance_reviews_user_idx ON performance_reviews(user_id);

CREATE TABLE match_review_likes (
    review_id UUID NOT NULL REFERENCES match_reviews(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    PRIMARY KEY (review_id, user_id)
);

CREATE INDEX match_review_likes_user_idx ON match_review_likes(user_id);

CREATE TABLE performance_review_likes (
    review_id UUID NOT NULL REFERENCES performance_reviews(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    PRIMARY KEY (review_id, user_id)
);

CREATE INDEX performance_review_likes_user_idx ON performance_review_likes(user_id);

CREATE TABLE match_review_comments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    review_id UUID NOT NULL REFERENCES match_reviews(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    body TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT match_review_comments_body_not_blank CHECK (btrim(body) <> '')
);

CREATE INDEX match_review_comments_review_idx ON match_review_comments(review_id);
CREATE INDEX match_review_comments_user_idx ON match_review_comments(user_id);

CREATE TABLE performance_review_comments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    review_id UUID NOT NULL REFERENCES performance_reviews(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    body TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT performance_review_comments_body_not_blank CHECK (btrim(body) <> '')
);

CREATE INDEX performance_review_comments_review_idx ON performance_review_comments(review_id);
CREATE INDEX performance_review_comments_user_idx ON performance_review_comments(user_id);

CREATE TABLE match_review_comment_likes (
    comment_id UUID NOT NULL REFERENCES match_review_comments(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    PRIMARY KEY (comment_id, user_id)
);

CREATE INDEX match_review_comment_likes_user_idx ON match_review_comment_likes(user_id);
