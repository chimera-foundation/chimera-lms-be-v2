package domain

import (
	"context"

	"github.com/google/uuid"
)

type ContentRepository interface {
	Create(ctx context.Context, content *Content) error
	GetByLessonID(ctx context.Context, lessonID uuid.UUID) ([]*Content, error)
}
