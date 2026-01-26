package domain

import (
	"context"

	"github.com/google/uuid"
)

type SectionRepository interface {
	Create(ctx context.Context, section *Section) error
	GetSectionIDsByUserID(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Section, error)
}
