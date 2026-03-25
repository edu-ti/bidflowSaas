package auth

import (
	"context"
	"database/sql"

	authqueries "lastsaas/core/auth/queries"

	"github.com/google/uuid"
)

type Service struct {
	db *sql.DB
	q  *authqueries.Queries
}

func NewService(db *sql.DB) *Service {
	return &Service{
		db: db,
		q:  authqueries.New(db),
	}
}

func (s *Service) CreateUser(ctx context.Context, email, passwordHash, firstName, lastName string, createdBy uuid.NullUUID) (authqueries.User, error) {
	return s.q.CreateUser(ctx, authqueries.CreateUserParams{
		Email:        email,
		PasswordHash: passwordHash,
		FirstName:    sql.NullString{String: firstName, Valid: firstName != ""},
		LastName:     sql.NullString{String: lastName, Valid: lastName != ""},
		Status:       "active",
		CreatedBy:    createdBy,
	})
}

func (s *Service) GetUserByEmail(ctx context.Context, email string) (authqueries.User, error) {
	return s.q.GetUserByEmail(ctx, email)
}

func (s *Service) GetUser(ctx context.Context, id uuid.UUID) (authqueries.User, error) {
	return s.q.GetUser(ctx, id)
}

func (s *Service) CreateTenantMember(ctx context.Context, tenantId, userId uuid.UUID, role string, createdBy uuid.NullUUID) (authqueries.TenantMember, error) {
	return s.q.CreateTenantMember(ctx, authqueries.CreateTenantMemberParams{
		TenantID:  tenantId,
		UserID:    userId,
		Role:      role,
		CreatedBy: createdBy,
	})
}

func (s *Service) GetUserTenants(ctx context.Context, userId uuid.UUID) ([]authqueries.GetUserTenantsRow, error) {
	return s.q.GetUserTenants(ctx, userId)
}
