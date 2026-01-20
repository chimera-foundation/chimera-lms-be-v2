package domain

import (
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/shared"
	u "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/user/domain"
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

	Instructor u.User
	Title string
	Description string
	Status CourseStatus	
	Price int64 
	Modules []m.Module
}