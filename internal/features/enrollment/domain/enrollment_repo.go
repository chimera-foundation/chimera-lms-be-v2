package domain

import (
	"context"

	"github.com/google/uuid"
)

type EnrollmentRepository interface {
	GetActiveSectionIDsByUserID(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error)
}