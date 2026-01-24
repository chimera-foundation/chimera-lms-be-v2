package domain

import (
    "context"
    "github.com/google/uuid"
)

type SectionMemberRepository interface {
    GetSectionIDsByUserID(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error)
}