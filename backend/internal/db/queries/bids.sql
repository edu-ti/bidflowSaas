-- name: CreateBid :one
INSERT INTO bids (
    tenant_id, title, description, published_date, status, portal
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetBid :one
SELECT * FROM bids
WHERE id = $1 AND tenant_id = $2;

-- name: ListBids :many
SELECT * FROM bids
WHERE tenant_id = $1
ORDER BY created_at DESC;
