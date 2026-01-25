package domain

import (
    "context"
    "github.com/google/uuid"
)

type CohortRepository interface {
    GetIDsByUserID(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error)
    GetByID(ctx context.Context, id uuid.UUID) (*Cohort, error)
}