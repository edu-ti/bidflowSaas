-- name: ListPlans :many
SELECT * FROM plans
WHERE deleted_at IS NULL
ORDER BY name;

-- name: GetTenantActiveModules :many
SELECT module_name FROM tenant_modules
WHERE tenant_id = $1 AND is_active = true AND deleted_at IS NULL;

-- name: CreateSubscription :one
INSERT INTO subscriptions (tenant_id, plan_id, status, start_date, created_by)
VALUES ($1, $2, $3, NOW(), $4)
RETURNING *;
