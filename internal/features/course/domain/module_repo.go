package domain

import (
	"context"

	"github.com/google/uuid"
)

type ModuleRepository interface {
	Create(ctx context.Context, module *Module) error
	GetByID(ctx context.Context, id uuid.UUID) (*Module, error)
}
