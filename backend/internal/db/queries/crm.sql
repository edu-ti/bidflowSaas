-- name: CreateCRMFunnel :one
INSERT INTO crm_funnels (
    tenant_id, bid_id, stage, notes
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetCRMFunnel :one
SELECT * FROM crm_funnels
WHERE id = $1 AND tenant_id = $2;

-- name: ListCRMFunnelsByBid :many
SELECT * FROM crm_funnels
WHERE tenant_id = $1 AND bid_id = $2
ORDER BY created_at DESC;
