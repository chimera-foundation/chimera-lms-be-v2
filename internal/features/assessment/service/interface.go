package service

import (
	"context"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/assessment/delivery/dto"
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/assessment/domain"
	"github.com/google/uuid"
)

// AssessmentService defines the interface for assessment business logic
type AssessmentService interface {
	GetStudentAssessments(ctx context.Context, userID uuid.UUID, filter domain.StudentAssessmentFilter) (*dto.StudentAssessmentsResponse, error)
}
