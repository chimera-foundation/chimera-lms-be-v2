package domain

import (
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/shared"
	"github.com/google/uuid"
)

type Section struct {
	shared.Base

	CohortID uuid.UUID

	Name string
	RoomCode string
	Capacity int
}