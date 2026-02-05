package service

import (
	"context"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/assessment/delivery/dto"
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/assessment/domain"
	"github.com/google/uuid"
)

type assessmentService struct {
	repo domain.AssessmentRepo
}

// NewAssessmentService creates a new assessment service
func NewAssessmentService(repo domain.AssessmentRepo) AssessmentService {
	return &assessmentService{
		repo: repo,
	}
}

func (s *assessmentService) GetStudentAssessments(ctx context.Context, userID uuid.UUID, filter domain.StudentAssessmentFilter) (*dto.StudentAssessmentsResponse, error) {
	// Get summary
	summary, err := s.repo.GetStudentAssessmentSummary(ctx, userID, filter)
	if err != nil {
		return nil, err
	}

	// Get assessment list
	items, err := s.repo.GetStudentAssessments(ctx, userID, filter)
	if err != nil {
		return nil, err
	}

	// Map to DTOs
	assessments := make([]dto.AssessmentItem, len(items))
	for i, item := range items {
		assessments[i] = dto.AssessmentItem{
			ID:            item.ID,
			Subject:       item.Subject,
			Title:         item.Title,
			AttachmentURL: item.AttachmentURL,
			Status:        string(item.Status),
			Type:          string(item.Type),
			SubType:       string(item.SubType),
			DueDate:       item.DueDate,
		}
	}

	return &dto.StudentAssessmentsResponse{
		Summary: dto.AssessmentSummary{
			Pending:   summary.Pending,
			Submitted: summary.Submitted,
			Done:      summary.Done,
			Overdue:   summary.Overdue,
		},
		Assessments: assessments,
	}, nil
}
