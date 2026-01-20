package domain

import (
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/shared"
	"github.com/google/uuid"
)

type Program struct {
	shared.Base

	OrganizationID uuid.UUID

	Name string
	Description string
}