package service

import (
	"context"
	"io"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/attachment/delivery/dto"
	"github.com/google/uuid"
)

type UploadRequest struct {
	OrgID        uuid.UUID
	UploadedBy   uuid.UUID
	AssessmentID *uuid.UUID
	SubmissionID *uuid.UUID
	FileName     string
	FileSize     int64
	MIMEType     string
	File         io.Reader
}

type AttachmentService interface {
	Upload(ctx context.Context, req UploadRequest) (*dto.AttachmentResponse, error)
	Delete(ctx context.Context, userID uuid.UUID, attachmentID uuid.UUID) error
	ListByAssessment(ctx context.Context, assessmentID uuid.UUID) ([]dto.AttachmentResponse, error)
	ListBySubmission(ctx context.Context, submissionID uuid.UUID) ([]dto.AttachmentResponse, error)
}
