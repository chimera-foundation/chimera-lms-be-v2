package domain

import (
	"context"

	"github.com/google/uuid"
)

type CourseRepository interface {
	Create(ctx context.Context, course *Course) error
	GetByID(ctx context.Context, id uuid.UUID) (*Course, error)
}
