package domain

import (
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/shared"
	"github.com/google/uuid"
)

type Module struct {
	shared.Base

	CourseID uuid.UUID

	Title string
	OrderIndex int
}