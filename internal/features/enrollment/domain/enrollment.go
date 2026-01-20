package domain

import (
	"time"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/shared"
	"github.com/google/uuid"
)

type EnrollmentStatus int

const (	
	Active EnrollmentStatus = iota
	Completed
	Dropped
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