-- name: CreateTenant :one
INSERT INTO tenants (name, slug, status, created_by)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetTenant :one
SELECT * FROM tenants
WHERE id = $1 AND deleted_at IS NULL LIMIT 1;

-- name: ListTenants :many
SELECT * FROM tenants
WHERE deleted_at IS NULL
ORDER BY name;

-- name: UpdateTenant :one
UPDATE tenants
SET name = $2, slug = $3, status = $4, updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: SoftDeleteTenant :exec
UPDATE tenants
SET deleted_at = NOW()
WHERE id = $1;
