package domain

import (
	"context"

	"github.com/google/uuid"
)

type OrganizationRepository interface {
	Create(ctx context.Context, org *Organization) error
	Update(ctx context.Context, org *Organization) error
	Delete(ctx context.Context, orgID uuid.UUID) error
	GetIDByUserID(ctx context.Context, userID uuid.UUID) (uuid.UUID, error)
}