package domain

import (
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/shared"
	s "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/subject/domain"
	c "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/course/domain"
)

type EducationLevel struct {
	shared.Base
	
	Name string
	Code string	
	Subjects []s.Subject
	Courses []c.Course
}