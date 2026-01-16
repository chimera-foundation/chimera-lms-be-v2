package ports

import (
	"context"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/modules/user/domain"
)

type Repository interface {
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
}

type UseCase interface {
	LoginByEmail(ctx context.Context, email, password string) (string, error)
	LoginByUsername(ctx context.Context, username, password string) (string, error)
}