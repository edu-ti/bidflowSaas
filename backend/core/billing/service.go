package billing

import (
	"context"
	"database/sql"

	billingqueries "lastsaas/core/billing/queries"

	"github.com/google/uuid"
)

type Service struct {
	db *sql.DB
	q  *billingqueries.Queries
}

func NewService(db *sql.DB) *Service {
	return &Service{
		db: db,
		q:  billingqueries.New(db),
	}
}

func (s *Service) TenantHasModule(ctx context.Context, tenantId uuid.UUID, moduleName string) (bool, error) {
	modules, err := s.q.GetTenantActiveModules(ctx, tenantId)
	if err != nil {
		return false, err
	}
	for _, mod := range modules {
		if mod == moduleName {
			return true, nil
		}
	}
	return false, nil
}

func (s *Service) ListPlans(ctx context.Context) ([]billingqueries.Plan, error) {
	return s.q.ListPlans(ctx)
}
