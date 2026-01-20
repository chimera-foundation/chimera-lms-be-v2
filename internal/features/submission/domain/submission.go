package domain

import (
	"time"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/shared"
	"github.com/google/uuid"
)

type Submission struct {
	shared.Base

	AssessmentID uuid.UUID
	UserID uuid.UUID
	EnrollmentID uuid.UUID

	FinalScore float32
	SubmittedAt time.Time
}