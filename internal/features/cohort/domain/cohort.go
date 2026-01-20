package domain

import (
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/shared"
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/section/domain"
)

type Cohort struct {
	shared.Base

	Name string
	Sections []domain.Section
}