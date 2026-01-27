package domain

import (
	"context"

	"github.com/google/uuid"
)

type ProgramRepository interface {
	Create(ctx context.Context, program *Program) error
	GetByID(ctx context.Context, id uuid.UUID) (*Program, error)
}
