-- name: CreateAIInsight :one
INSERT INTO ai_insights (
    tenant_id, bid_id, embedding, insight_text
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetAIInsight :one
SELECT * FROM ai_insights
WHERE id = $1 AND tenant_id = $2;

-- name: SearchSimilarInsights :many
SELECT *, 1 - (embedding <=> $3::vector) AS similarity
FROM ai_insights
WHERE tenant_id = $1 AND bid_id = $2
ORDER BY embedding <=> $3::vector
LIMIT 5;
