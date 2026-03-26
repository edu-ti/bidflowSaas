package crm

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"

	"lastsaas/api/middleware"
	crmrepo "lastsaas/modules/crm/repository"
	crmqueries "lastsaas/modules/crm/queries"
)

type Service struct {
	repo crmrepo.Repository
}

func NewService(db *sql.DB) *Service {
	return &Service{
		repo: crmrepo.NewRepository(db),
	}
}

// CreateLeadInput represents the required fields to create a Lead
type CreateLeadInput struct {
	Title      string
	CustomerID uuid.NullUUID
	Status     string
	Value      string
	CreatedBy  uuid.NullUUID
}

func (s *Service) CreateLead(ctx context.Context, input CreateLeadInput) (*crmqueries.Lead, error) {
	// 1. Extract tenant
	tenantIDStr := middleware.GetTenantID(ctx)
	if tenantIDStr == "" {
		return nil, errors.New("unauthorized: missing tenant context")
	}

	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		return nil, errors.New("invalid tenant ID format")
	}

	// 2. Business Validations
	if input.Title == "" {
		return nil, errors.New("validation failed: lead title is required")
	}
	if input.Status == "" {
		input.Status = "NEW" // default status
	}

	// Example: Duplicate Lead Title validation
	exists, err := s.repo.GetLeadByTitleAndTenant(ctx, input.Title, tenantID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("validation failed: lead with this title already exists")
	}

	// 3. Transaction support for multi-step operations
	// Example of using a transaction if we had multiple steps
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	txRepo := s.repo.WithTx(tx)

	// Create Lead
	lead, err := txRepo.CreateLead(ctx, crmqueries.CreateLeadParams{
		TenantID:   tenantID,
		Title:      input.Title,
		CustomerID: input.CustomerID,
		Status:     input.Status,
		CreatedBy:  input.CreatedBy,
		// Assuming Value would be parsed if it existed in the SQL schema
	})
	if err != nil {
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &lead, nil
}

func (s *Service) ListLeads(ctx context.Context) ([]crmqueries.Lead, error) {
	tenantIDStr := middleware.GetTenantID(ctx)
	if tenantIDStr == "" {
		return nil, errors.New("unauthorized: missing tenant context")
	}

	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		return nil, errors.New("invalid tenant ID format")
	}

	return s.repo.ListLeads(ctx, tenantID)
}

// The following methods will need to be refactored similar to CreateLead 
// to complete the repository pattern implementation for full coverage.
func (s *Service) CreateCustomer(ctx context.Context, name, email, phone, documentId string, createdBy uuid.NullUUID) (crmqueries.Customer, error) {
	return crmqueries.Customer{}, errors.New("not implemented: pending refactor to repository")
}

func (s *Service) ListCustomers(ctx context.Context) ([]crmqueries.Customer, error) {
	return nil, errors.New("not implemented: pending refactor to repository")
}

func (s *Service) CreateOpportunity(ctx context.Context, leadId, customerId uuid.UUID, stage string, expectedCloseDate sql.NullTime, createdBy uuid.NullUUID) (crmqueries.Opportunity, error) {
	return crmqueries.Opportunity{}, errors.New("not implemented: pending refactor to repository")
}

func (s *Service) ListOpportunities(ctx context.Context) ([]crmqueries.Opportunity, error) {
	return nil, errors.New("not implemented: pending refactor to repository")
}

