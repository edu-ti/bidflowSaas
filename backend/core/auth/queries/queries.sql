-- name: CreateUser :one
INSERT INTO users (email, password_hash, first_name, last_name, status, created_by)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 AND deleted_at IS NULL LIMIT 1;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 AND deleted_at IS NULL LIMIT 1;

-- name: CreateTenantMember :one
INSERT INTO tenant_members (tenant_id, user_id, role, created_by)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetUserTenants :many
SELECT t.*, tm.role 
FROM tenants t
JOIN tenant_members tm ON t.id = tm.tenant_id
WHERE tm.user_id = $1 AND t.deleted_at IS NULL AND tm.deleted_at IS NULL;
