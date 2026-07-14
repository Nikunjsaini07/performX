-- Trending scores: precomputed by the background worker every 30 minutes.
-- The API only reads from this table; it never computes trending on the fly.

CREATE TABLE trending_scores (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    entity_type TEXT NOT NULL,
    entity_id UUID NOT NULL,
    time_window TEXT NOT NULL,
    score NUMERIC(12,2) NOT NULL DEFAULT 0,
    rank INTEGER NOT NULL DEFAULT 0,
    calculated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT trending_scores_entity_window_unique UNIQUE (entity_type, entity_id, time_window),
    CONSTRAINT trending_scores_entity_type_valid CHECK (
        entity_type IN ('performance', 'player', 'match', 'review')
    ),
    CONSTRAINT trending_scores_window_valid CHECK (
        time_window IN ('today', 'week', 'month')
    ),
    CONSTRAINT trending_scores_score_non_negative CHECK (score >= 0),
    CONSTRAINT trending_scores_rank_positive CHECK (rank > 0)
);

-- Primary read path: type + window ordered by rank
CREATE INDEX trending_scores_type_window_rank_idx
    ON trending_scores(entity_type, time_window, rank);

-- Cleanup path: find stale entries
CREATE INDEX trending_scores_calculated_at_idx
    ON trending_scores(calculated_at);
