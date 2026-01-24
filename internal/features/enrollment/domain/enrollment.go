package domain

import (
	"time"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/shared"
	"github.com/google/uuid"
)

type EnrollmentStatus string

const (	
	Active EnrollmentStatus = "active"
	Completed EnrollmentStatus = "completed"
	Dropped EnrollmentStatus = "dropped"
)

type Enrollment struct {
	shared.Base

	UserID uuid.UUID
	CourseID uuid.UUID
	SectionID uuid.UUID
	AcademicPeriodID uuid.UUID

	Status EnrollmentStatus
	EnrolledAt time.Time
}