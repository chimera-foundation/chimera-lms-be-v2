package domain

import (
	"time"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/shared"
	"github.com/google/uuid"
)

type ProgressTracker struct {
	shared.Base

	EnrollmentID uuid.UUID
	ContentID uuid.UUID

	IsCompleted bool
	UpdatedAt time.Time
}