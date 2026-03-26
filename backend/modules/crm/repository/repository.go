package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	crmqueries "lastsaas/modules/crm/queries"
)

type Repository interface {
	CreateLead(ctx context.Context, arg crmqueries.CreateLeadParams) (crmqueries.Lead, error)
	ListLeads(ctx context.Context, tenantID uuid.UUID) ([]crmqueries.Lead, error)
	GetLeadByTitleAndTenant(ctx context.Context, title string, tenantID uuid.UUID) (bool, error)
	WithTx(tx *sql.Tx) Repository
	BeginTx(ctx context.Context) (*sql.Tx, error)
}

type repository struct {
	db *sql.DB
	q  *crmqueries.Queries
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
		q:  crmqueries.New(db),
	}
}

func (r *repository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, nil)
}

func (r *repository) WithTx(tx *sql.Tx) Repository {
	return &repository{
		db: r.db, // Keep original DB just in case, but queries will use tx
		q:  r.q.WithTx(tx),
	}
}

func (r *repository) CreateLead(ctx context.Context, arg crmqueries.CreateLeadParams) (crmqueries.Lead, error) {
	return r.q.CreateLead(ctx, arg)
}

func (r *repository) ListLeads(ctx context.Context, tenantID uuid.UUID) ([]crmqueries.Lead, error) {
	return r.q.ListLeads(ctx, tenantID)
}

func (r *repository) GetLeadByTitleAndTenant(ctx context.Context, title string, tenantID uuid.UUID) (bool, error) {
	var exists bool
	// Standard validation query to check for duplicate
	err := r.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM crm_leads WHERE title = $1 AND tenant_id = $2)", title, tenantID).Scan(&exists)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return false, err
	}
	return exists, nil
}
