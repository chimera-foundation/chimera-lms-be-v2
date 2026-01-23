package service

import (
	"context"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/user/domain"
	"github.com/google/uuid"
)

type Auth interface {
    Register(ctx context.Context, email, password, firstName, lastName string, orgID uuid.UUID) (*domain.User, error)
    Login(ctx context.Context, email, password string) (string, error)
	Logout(ctx context.Context, token string) (error)
}