package domain

import (
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/shared"
	"github.com/google/uuid"
)

type CohortMember struct {
	shared.Base

	CohortID uuid.UUID
	UserID   uuid.UUID
}
