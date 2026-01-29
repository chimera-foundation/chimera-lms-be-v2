package auth

import (
	"context"

	"github.com/google/uuid"
)


type contextKey string
const (
	UserIDKey contextKey = "userID"
	OrgIDKey contextKey = "orgID"
)

func GetUserID(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(UserIDKey).(uuid.UUID)
	return userID, ok
}

func GetOrgID(ctx context.Context) (uuid.UUID, bool) {
	orgID, ok := ctx.Value(OrgIDKey).(uuid.UUID)
	return orgID, ok
}