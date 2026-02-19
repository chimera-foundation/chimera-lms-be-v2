package service

import (
	"context"
	"fmt"
	"time"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/attachment/delivery/dto"
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/attachment/domain"
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/shared/storage"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type attachmentService struct {
	repo    domain.AttachmentRepo
	storage storage.FileStorage
	log     *logrus.Logger
}

func NewAttachmentService(repo domain.AttachmentRepo, storage storage.FileStorage, log *logrus.Logger) AttachmentService {
	return &attachmentService{
		repo:    repo,
		storage: storage,
		log:     log,
	}
}

func (s *attachmentService) Upload(ctx context.Context, req UploadRequest) (*dto.AttachmentResponse, error) {
	attachment := &domain.Attachment{
		ID:             uuid.New(),
		OrganizationID: req.OrgID,
		UploadedBy:     req.UploadedBy,
		AssessmentID:   req.AssessmentID,
		SubmissionID:   req.SubmissionID,
		FileName:       req.FileName,
		FileSize:       req.FileSize,
		MIMEType:       req.MIMEType,
		CreatedAt:      time.Now(),
	}

	if err := attachment.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	if req.AssessmentID != nil {
		count, err := s.repo.CountByAssessment(ctx, *req.AssessmentID)
		if err != nil {
			return nil, fmt.Errorf("failed to count assessment attachments: %w", err)
		}
		if count >= domain.MaxAttachmentsPerEntity {
			return nil, fmt.Errorf("maximum number of attachments (%d) reached for this assessment", domain.MaxAttachmentsPerEntity)
		}
	}
	if req.SubmissionID != nil {
		count, err := s.repo.CountBySubmission(ctx, *req.SubmissionID)
		if err != nil {
			return nil, fmt.Errorf("failed to count submission attachments: %w", err)
		}
		if count >= domain.MaxAttachmentsPerEntity {
			return nil, fmt.Errorf("maximum number of attachments (%d) reached for this submission", domain.MaxAttachmentsPerEntity)
		}
	}

	storagePath := fmt.Sprintf("attachments/%s/%s", attachment.ID.String(), req.FileName)
	fileURL, err := s.storage.Upload(ctx, storagePath, req.File)
	if err != nil {
		s.log.WithError(err).Error("failed to upload file to storage")
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}
	attachment.FileURL = fileURL

	if err := s.repo.Create(ctx, attachment); err != nil {
		_ = s.storage.Delete(ctx, storagePath)
		return nil, fmt.Errorf("failed to save attachment: %w", err)
	}

	return toDTO(attachment), nil
}

func (s *attachmentService) Delete(ctx context.Context, userID uuid.UUID, attachmentID uuid.UUID) error {
	attachment, err := s.repo.GetByID(ctx, attachmentID)
	if err != nil {
		return fmt.Errorf("failed to get attachment: %w", err)
	}
	if attachment == nil {
		return fmt.Errorf("attachment not found")
	}

	if attachment.UploadedBy != userID {
		return fmt.Errorf("you are not authorized to delete this attachment")
	}

	if err := s.repo.SoftDelete(ctx, attachmentID); err != nil {
		return fmt.Errorf("failed to delete attachment: %w", err)
	}

	storagePath := fmt.Sprintf("attachments/%s/%s", attachment.ID.String(), attachment.FileName)
	if delErr := s.storage.Delete(ctx, storagePath); delErr != nil {
		s.log.WithError(delErr).WithField("attachment_id", attachmentID).Warn("failed to delete file from storage")
	}

	return nil
}

func (s *attachmentService) ListByAssessment(ctx context.Context, assessmentID uuid.UUID) ([]dto.AttachmentResponse, error) {
	attachments, err := s.repo.ListByAssessment(ctx, assessmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to list assessment attachments: %w", err)
	}
	return toDTOList(attachments), nil
}

func (s *attachmentService) ListBySubmission(ctx context.Context, submissionID uuid.UUID) ([]dto.AttachmentResponse, error) {
	attachments, err := s.repo.ListBySubmission(ctx, submissionID)
	if err != nil {
		return nil, fmt.Errorf("failed to list submission attachments: %w", err)
	}
	return toDTOList(attachments), nil
}

// --- helpers ---

func toDTO(a *domain.Attachment) *dto.AttachmentResponse {
	return &dto.AttachmentResponse{
		ID:        a.ID,
		FileName:  a.FileName,
		FileURL:   a.FileURL,
		FileSize:  a.FileSize,
		MIMEType:  a.MIMEType,
		CreatedAt: a.CreatedAt,
	}
}

func toDTOList(attachments []domain.Attachment) []dto.AttachmentResponse {
	result := make([]dto.AttachmentResponse, len(attachments))
	for i, a := range attachments {
		result[i] = dto.AttachmentResponse{
			ID:        a.ID,
			FileName:  a.FileName,
			FileURL:   a.FileURL,
			FileSize:  a.FileSize,
			MIMEType:  a.MIMEType,
			CreatedAt: a.CreatedAt,
		}
	}
	return result
}
