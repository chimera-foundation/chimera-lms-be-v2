package domain

import (
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/shared"
	"github.com/google/uuid"
)

type Cohort struct {
	shared.Base

	OrganizationID uuid.UUID
	AcademicPeriodID uuid.UUID
	EducationLevelID uuid.UUID

	Name string
}