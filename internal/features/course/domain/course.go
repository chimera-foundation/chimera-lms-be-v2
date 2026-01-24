package domain

import (
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/shared"
	"github.com/google/uuid"
)

type CourseStatus string

const (
	Draft CourseStatus = "draft"
	Published CourseStatus = "published"
	Archived CourseStatus = "archived"
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