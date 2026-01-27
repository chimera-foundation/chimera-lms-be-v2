package domain

import (
	"context"

	"github.com/google/uuid"
)

type RoleRepository interface {
	Create(ctx context.Context, role *Role) error
	GetByName(ctx context.Context, name string) (*Role, error)
	AssignRoleToUser(ctx context.Context, userID, roleID uuid.UUID) error
}
