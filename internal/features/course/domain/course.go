package domain

import (
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/shared"
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/user/domain"
)

type CourseStatus int

const (
	Draft CourseStatus = iota
	Published
	Archived
)

type Course struct {
	shared.Base

	Instructor domain.User
	Title string
	Description string
	Status CourseStatus	
	Price int64 
}