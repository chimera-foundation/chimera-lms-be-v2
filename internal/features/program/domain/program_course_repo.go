package domain

import (
	"context"
)

type ProgramCourseRepository interface {
	Create(ctx context.Context, pc *ProgramCourse) error
}
