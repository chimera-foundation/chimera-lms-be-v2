package domain

import (
	"context"

	"github.com/google/uuid"
)

type SubjectRepository interface {
	Create(ctx context.Context, subject *Subject) error
	GetByID(ctx context.Context, id uuid.UUID) (*Subject, error)
}
