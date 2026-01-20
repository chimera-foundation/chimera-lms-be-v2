package domain

import (
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/shared"
	l "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/lesson/domain"
)

type Module struct {
	shared.Base

	Title string
	OrderIndex int

	Lessons []l.Lesson
}