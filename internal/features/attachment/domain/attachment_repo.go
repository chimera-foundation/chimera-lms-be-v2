package domain

import (
	"context"

	"github.com/google/uuid"
)

type AttachmentRepo interface {
	Create(ctx context.Context, attachment *Attachment) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*Attachment, error)
	ListByAssessment(ctx context.Context, assessmentID uuid.UUID) ([]Attachment, error)
	ListBySubmission(ctx context.Context, submissionID uuid.UUID) ([]Attachment, error)
	CountByAssessment(ctx context.Context, assessmentID uuid.UUID) (int, error)
	CountBySubmission(ctx context.Context, submissionID uuid.UUID) (int, error)
}
