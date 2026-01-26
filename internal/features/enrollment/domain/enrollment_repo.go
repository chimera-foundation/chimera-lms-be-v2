package domain

import (
	"context"

	"github.com/google/uuid"
)

type EnrollmentRepository interface {
	Create(ctx context.Context, enrollment *Enrollment) error
	GetActiveSectionIDsByUserID(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error)
}
