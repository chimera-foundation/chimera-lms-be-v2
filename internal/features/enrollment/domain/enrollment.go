package domain

import (
	"time"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/shared"
	sub "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/submission/domain"
)

type EnrollmentStatus int

const (	
	Active EnrollmentStatus = iota
	Completed
	Dropped
)

type Enrollment struct {
	shared.Base

	Status EnrollmentStatus
	EnrolledAt time.Time
	Submissions []sub.Submission
}