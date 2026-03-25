-- name: CreateLog :exec
INSERT INTO system_logs (tenant_id, user_id, action)
VALUES ($1, $2, $3);

-- name: ListTenantLogs :many
SELECT * FROM system_logs
WHERE tenant_id = $1
ORDER BY timestamp DESC
LIMIT $2 OFFSET $3;
