package domain

import (
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/shared"
	"github.com/google/uuid"
)

type CourseStatus int

const (
	Draft CourseStatus = iota
	Published
	Archived
)

type Course struct {
	shared.Base

	OrganizationID uuid.UUID
	InstructorID uuid.UUID
	SubjectID uuid.UUID
	EducationLevelID uuid.UUID

	Title string
	Description string
	Status CourseStatus	
	Price int64 
	GradeLevel int
	Credits int // for university
}