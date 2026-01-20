package domain

import (
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/shared"
	"github.com/google/uuid"
)

type Lesson struct {
	shared.Base

	ModuleID uuid.UUID

	Title string
	OrderIndex int
}