package domain

import (
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/shared"
	"github.com/google/uuid"
)

type Subject struct {
	shared.Base

	OrganizationID uuid.UUID
	EducationLevelID uuid.UUID

	Name string
	Code string
}