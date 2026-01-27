package domain

import (
	"context"

	"github.com/google/uuid"
)

type AcademicPeriodRepository interface {
	Create(ctx context.Context, period *AcademicPeriod, orgID uuid.UUID) error
	GetActiveByOrganizationID(ctx context.Context, orgID uuid.UUID) (*AcademicPeriod, error)
}
