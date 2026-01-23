package auth

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type TokenProvider interface {
	GenerateToken(userID uuid.UUID, orgID uuid.UUID) (string, error)
	ValidateToken(token string) (uuid.UUID, error)
	BlacklistToken(ctx context.Context, token string, expiration time.Duration) error
    IsBlacklisted(ctx context.Context, token string) (bool, error)
}