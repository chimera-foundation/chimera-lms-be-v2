package domain

import (
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/shared"
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/course/domain"
)

type Program struct {
	shared.Base

	Name string
	Description string
	Courses []domain.Course
}