-- name: GetTrendingScores :many
SELECT
    ts.id,
    ts.entity_type,
    ts.entity_id,
    ts.time_window,
    ts.score,
    ts.rank,
    ts.calculated_at
FROM trending_scores ts
WHERE ts.entity_type = $1
  AND ts.time_window = $2
ORDER BY ts.rank ASC
LIMIT $3;

-- name: UpsertTrendingScore :exec
INSERT INTO trending_scores (entity_type, entity_id, time_window, score, rank, calculated_at)
VALUES ($1, $2, $3, $4, $5, now())
ON CONFLICT (entity_type, entity_id, time_window)
DO UPDATE SET
    score = EXCLUDED.score,
    rank = EXCLUDED.rank,
    calculated_at = now();

-- name: DeleteTrendingScoresByTypeAndWindow :exec
DELETE FROM trending_scores
WHERE entity_type = $1
  AND time_window = $2;

-- name: CountTrendingScores :one
SELECT COUNT(*) FROM trending_scores
WHERE entity_type = $1
  AND time_window = $2;
