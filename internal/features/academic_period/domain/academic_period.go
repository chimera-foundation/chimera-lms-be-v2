package domain

import (
	"time"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/shared"
	co "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/cohort/domain"
	e "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/enrollment/domain"
)

type AcademicPeriod struct {
	shared.Base

	Name string // e.g., "2025/2026 ganjil"
	StartDate time.Time
	EndDate time.Time
	Cohorts []co.Cohort
	Enrollments []e.Enrollment

	IsActive bool
}