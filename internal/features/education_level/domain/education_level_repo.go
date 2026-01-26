package domain

import (
	"context"

	"github.com/google/uuid"
)

type EducationLevelRepository interface {
	Create(ctx context.Context, level *EducationLevel) error
	GetByID(ctx context.Context, id uuid.UUID) (*EducationLevel, error)
}
