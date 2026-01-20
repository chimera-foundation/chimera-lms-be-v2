package domain

import (
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/shared"
)

type Cohort struct {
	shared.Base

	Name string
}