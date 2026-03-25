-- name: CreateEdital :one
INSERT INTO editais (tenant_id, number, agency, object_description, opening_date, estimated_value, status, created_by)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: ListEditais :many
SELECT * FROM editais
WHERE tenant_id = $1 AND deleted_at IS NULL
ORDER BY opening_date DESC;

-- name: CreateProposta :one
INSERT INTO propostas (tenant_id, edital_id, submitted_value, submission_date, status, created_by)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: ListPropostas :many
SELECT * FROM propostas
WHERE tenant_id = $1 AND edital_id = $2 AND deleted_at IS NULL
ORDER BY submission_date DESC;

-- name: RegisterResultado :one
INSERT INTO resultados (tenant_id, edital_id, proposta_id, won, notes, created_by)
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (tenant_id, edital_id) 
DO UPDATE SET 
    proposta_id = EXCLUDED.proposta_id, 
    won = EXCLUDED.won, 
    notes = EXCLUDED.notes, 
    updated_at = NOW()
RETURNING *;

-- name: GetResultado :one
SELECT * FROM resultados
WHERE tenant_id = $1 AND edital_id = $2 AND deleted_at IS NULL LIMIT 1;
