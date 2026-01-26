package domain

import (
	"context"

	"github.com/google/uuid"
)

type LessonRepository interface {
	Create(ctx context.Context, lesson *Lesson) error
	GetByID(ctx context.Context, id uuid.UUID) (*Lesson, error)
}
