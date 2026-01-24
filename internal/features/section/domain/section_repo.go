package domain

import (
    "context"
    "github.com/google/uuid"
)

type SectionRepository interface {
    GetSectionIDsByUserID(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error)
    GetByID(ctx context.Context, id uuid.UUID) (*Section, error)
}