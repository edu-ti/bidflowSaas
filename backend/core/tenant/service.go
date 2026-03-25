package tenant

import (
	"context"
	"database/sql"

	tenantqueries "lastsaas/core/tenant/queries"

	"github.com/google/uuid"
)

type Service struct {
	db *sql.DB
	q  *tenantqueries.Queries
}

func NewService(db *sql.DB) *Service {
	return &Service{
		db: db,
		q:  tenantqueries.New(db),
	}
}

func (s *Service) CreateTenant(ctx context.Context, name, slug string, createdBy uuid.NullUUID) (tenantqueries.Tenant, error) {
	return s.q.CreateTenant(ctx, tenantqueries.CreateTenantParams{
		Name:      name,
		Slug:      slug,
		Status:    "active",
		CreatedBy: createdBy,
	})
}

func (s *Service) GetTenant(ctx context.Context, id uuid.UUID) (tenantqueries.Tenant, error) {
	return s.q.GetTenant(ctx, id)
}

func (s *Service) ListTenants(ctx context.Context) ([]tenantqueries.Tenant, error) {
	return s.q.ListTenants(ctx)
}

func (s *Service) UpdateTenant(ctx context.Context, id uuid.UUID, name, slug, status string) (tenantqueries.Tenant, error) {
	return s.q.UpdateTenant(ctx, tenantqueries.UpdateTenantParams{
		ID:     id,
		Name:   name,
		Slug:   slug,
		Status: status,
	})
}

func (s *Service) SoftDeleteTenant(ctx context.Context, id uuid.UUID) error {
	return s.q.SoftDeleteTenant(ctx, id)
}
