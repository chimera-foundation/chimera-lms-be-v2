package domain

import (
	"context"

	"github.com/google/uuid"
)

type AssessmentRepo interface {
	Create(ctx context.Context, assessment *Assessment) error
	GetStudentAssessments(ctx context.Context, userID uuid.UUID, filter StudentAssessmentFilter) ([]StudentAssessmentItem, error)
	GetStudentAssessmentSummary(ctx context.Context, userID uuid.UUID, filter StudentAssessmentFilter) (*StudentAssessmentSummary, error)
}
