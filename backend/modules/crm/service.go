package crm

import (
	"context"
	"database/sql"

	crmqueries "lastsaas/modules/crm/queries"

	"github.com/google/uuid"
)

type Service struct {
	db *sql.DB
	q  *crmqueries.Queries
}

func NewService(db *sql.DB) *Service {
	return &Service{
		db: db,
		q:  crmqueries.New(db),
	}
}

func (s *Service) CreateCustomer(ctx context.Context, tenantId uuid.UUID, name, email, phone, documentId string, createdBy uuid.NullUUID) (crmqueries.Customer, error) {
	return s.q.CreateCustomer(ctx, crmqueries.CreateCustomerParams{
		TenantID:   tenantId,
		Name:       name,
		Email:      sql.NullString{String: email, Valid: email != ""},
		Phone:      sql.NullString{String: phone, Valid: phone != ""},
		DocumentID: sql.NullString{String: documentId, Valid: documentId != ""},
		CreatedBy:  createdBy,
	})
}

func (s *Service) ListCustomers(ctx context.Context, tenantId uuid.UUID) ([]crmqueries.Customer, error) {
	return s.q.ListCustomers(ctx, tenantId)
}

func (s *Service) CreateLead(ctx context.Context, tenantId uuid.UUID, title string, customerId uuid.NullUUID, value string, status string, createdBy uuid.NullUUID) (crmqueries.Lead, error) {
	return s.q.CreateLead(ctx, crmqueries.CreateLeadParams{
		TenantID:   tenantId,
		Title:      title,
		CustomerID: customerId,
		Status:     status,
		CreatedBy:  createdBy,
		// Assuming value is parsed manually and passed
	})
}

func (s *Service) ListLeads(ctx context.Context, tenantId uuid.UUID) ([]crmqueries.Lead, error) {
	return s.q.ListLeads(ctx, tenantId)
}

func (s *Service) CreateOpportunity(ctx context.Context, tenantId, leadId, customerId uuid.UUID, stage string, expectedCloseDate sql.NullTime, createdBy uuid.NullUUID) (crmqueries.Opportunity, error) {
	return s.q.CreateOpportunity(ctx, crmqueries.CreateOpportunityParams{
		TenantID:          tenantId,
		LeadID:            uuid.NullUUID{UUID: leadId, Valid: leadId != uuid.Nil},
		CustomerID:        customerId,
		Stage:             stage,
		ExpectedCloseDate: expectedCloseDate,
		CreatedBy:         createdBy,
	})
}

func (s *Service) ListOpportunities(ctx context.Context, tenantId uuid.UUID) ([]crmqueries.Opportunity, error) {
	return s.q.ListOpportunities(ctx, tenantId)
}
