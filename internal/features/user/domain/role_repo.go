package domain

import (
	"context"
)

type RoleRepository interface {
	Create(ctx context.Context, role *Role) error
	GetByName(ctx context.Context, name string) (*Role, error)
}
