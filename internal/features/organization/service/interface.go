package service

import (
	"context"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/organization/domain"
	"github.com/google/uuid"
)

type CRUD interface {
	RegisterOrganization(ctx context.Context, name, slug, address string, isActive bool) (*domain.Organization, error)
	DeleteOrganization(ctx context.Context, orgID uuid.UUID) (error) 
	GetOrganization(ctx context.Context, orgID uuid.UUID) (*domain.Organization, error)
}