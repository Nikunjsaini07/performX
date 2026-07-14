-- Revert drop hierarchy
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

-- Revert stat_types
DROP INDEX IF EXISTS stat_types_category_display_order_idx;
ALTER TABLE stat_types DROP CONSTRAINT IF EXISTS stat_types_name_unique;
ALTER TABLE stat_types DROP CONSTRAINT IF EXISTS stat_types_short_name_unique;

ALTER TABLE stat_types ADD COLUMN sport_id UUID REFERENCES sports(id) ON DELETE RESTRICT;
-- In a real downgrade we'd have to populate sport_id before making it NOT NULL, but we'll leave it simple
ALTER TABLE stat_types ADD CONSTRAINT stat_types_sport_name_unique UNIQUE (sport_id, name);
ALTER TABLE stat_types ADD CONSTRAINT stat_types_sport_short_name_unique UNIQUE (sport_id, short_name);
CREATE INDEX stat_types_sport_category_display_order_idx ON stat_types (sport_id, category, display_order);

-- Revert players
ALTER TABLE players ADD COLUMN sport_id UUID REFERENCES sports(id) ON DELETE RESTRICT;
CREATE INDEX players_sport_id_idx ON players(sport_id);

-- Revert teams
ALTER TABLE teams ADD COLUMN sport_id UUID REFERENCES sports(id) ON DELETE RESTRICT;
CREATE INDEX teams_sport_idx ON teams(sport_id);

-- Revert matches
DROP INDEX IF EXISTS matches_utc_datetime_idx;
ALTER TABLE matches ADD COLUMN season_id UUID REFERENCES seasons(id) ON DELETE RESTRICT;
CREATE INDEX matches_season_datetime_idx ON matches(season_id, utc_datetime);

-- Revert Follows
CREATE TABLE user_follows (
    follower_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    following_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    PRIMARY KEY (follower_id, following_id),
    CONSTRAINT user_follows_no_self_follow CHECK (follower_id <> following_id)
);
CREATE INDEX user_follows_following_idx ON user_follows(following_id);

-- Revert Lists
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
