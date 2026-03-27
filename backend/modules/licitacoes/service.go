package licitacoes

import (
	"context"
	"database/sql"
	"time"

	"lastsaas/core/billing"
	licqueries "lastsaas/modules/licitacoes/queries"

	"github.com/google/uuid"
)

type Service struct {
	db *sql.DB
	q  *licqueries.Queries
}

func NewService(db *sql.DB) *Service {
	return &Service{
		db: db,
		q:  licqueries.New(db),
	}
}

func (s *Service) CreateEdital(ctx context.Context, tenantId uuid.UUID, number, agency, objectDesc string, openingDate time.Time, createdBy uuid.NullUUID) (licqueries.Editai, error) {
	result, err := s.q.CreateEdital(ctx, licqueries.CreateEditalParams{
		TenantID:          tenantId,
		Number:            number,
		Agency:            agency,
		ObjectDescription: objectDesc,
		OpeningDate:       openingDate,
		Status:            "open",
		CreatedBy:         createdBy,
	})
	if err == nil {
		_ = billing.IncrementUsage(ctx, s.db, tenantId, "editais")
	}
	return result, err
}

func (s *Service) ListEditais(ctx context.Context, tenantId uuid.UUID) ([]licqueries.Editai, error) {
	return s.q.ListEditais(ctx, tenantId)
}

func (s *Service) CreateProposta(ctx context.Context, tenantId, editalId uuid.UUID, value string, status string, createdBy uuid.NullUUID) (licqueries.Proposta, error) {
	result, err := s.q.CreateProposta(ctx, licqueries.CreatePropostaParams{
		TenantID: tenantId,
		EditalID: editalId,
		Status:   status,
		CreatedBy: createdBy,
	})
	if err == nil {
		_ = billing.IncrementUsage(ctx, s.db, tenantId, "propostas")
	}
	return result, err
}

func (s *Service) ListPropostas(ctx context.Context, tenantId, editalId uuid.UUID) ([]licqueries.Proposta, error) {
	return s.q.ListPropostas(ctx, licqueries.ListPropostasParams{
		TenantID: tenantId,
		EditalID: editalId,
	})
}

func (s *Service) RegisterResultado(ctx context.Context, tenantId, editalId uuid.UUID, propostaId uuid.NullUUID, won bool, notes string, createdBy uuid.NullUUID) (licqueries.Resultado, error) {
	return s.q.RegisterResultado(ctx, licqueries.RegisterResultadoParams{
		TenantID:   tenantId,
		EditalID:   editalId,
		PropostaID: propostaId,
		Won:        won,
		Notes:      sql.NullString{String: notes, Valid: notes != ""},
		CreatedBy:  createdBy,
	})
}

