package domain

import (
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/shared"
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/course/domain"
)

type Subject struct {
	shared.Base

	Name string
	Code string
	Courses []domain.Course
}