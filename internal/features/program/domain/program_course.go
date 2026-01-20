package domain

import (
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/shared"
	"github.com/google/uuid"
)

type ProgramCourse struct {
	shared.Base

	ProgramID uuid.UUID
	CourseID uuid.UUID
	OrderIndex int
}