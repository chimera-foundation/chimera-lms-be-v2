package domain

import (
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/shared"
	m "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/module/domain"
)

type CourseStatus int

const (
	Draft CourseStatus = iota
	Published
	Archived
)

type Course struct {
	shared.Base

	Title string
	Description string
	Status CourseStatus	
	Price int64 
	Modules []m.Module
	GradeLevel int
	Credits int // for university
}