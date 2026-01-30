package domain

import (
	"time"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/shared"
)

type AcademicPeriod struct {
	shared.Base

	Name string // e.g., "2025/2026 ganjil"
	StartDate time.Time
	EndDate time.Time

	IsActive bool
}

func NewAcademicPeriod(name string, startDate, endDate time.Time) *AcademicPeriod {
	return &AcademicPeriod{
		Name: name,
		StartDate: startDate,
		EndDate: endDate,
		IsActive: true,
	}
}