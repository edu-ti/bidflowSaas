-- name: CreateCustomer :one
INSERT INTO customers (tenant_id, name, email, phone, document_id, created_by)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: ListCustomers :many
SELECT * FROM customers
WHERE tenant_id = $1 AND deleted_at IS NULL
ORDER BY name;

-- name: CreateLead :one
INSERT INTO leads (tenant_id, title, customer_id, value, status, created_by)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: ListLeads :many
SELECT * FROM leads
WHERE tenant_id = $1 AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: CreateOpportunity :one
INSERT INTO opportunities (tenant_id, lead_id, customer_id, stage, amount, expected_close_date, created_by)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: ListOpportunities :many
SELECT * FROM opportunities
WHERE tenant_id = $1 AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: CreateActivity :one
INSERT INTO activities (tenant_id, related_to_type, related_to_id, type, description, due_date, created_by)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: ListActivities :many
SELECT * FROM activities
WHERE tenant_id = $1 AND deleted_at IS NULL
ORDER BY due_date ASC;
