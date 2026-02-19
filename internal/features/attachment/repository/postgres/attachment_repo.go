package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/attachment/domain"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type AttachmentRepoPostgres struct {
	db  *sql.DB
	log *logrus.Logger
}

func NewAttachmentRepoPostgres(db *sql.DB, log *logrus.Logger) domain.AttachmentRepo {
	return &AttachmentRepoPostgres{
		db:  db,
		log: log,
	}
}

func (r *AttachmentRepoPostgres) Create(ctx context.Context, attachment *domain.Attachment) error {
	query := `
		INSERT INTO attachments (
			id, organization_id, uploaded_by, assessment_id, submission_id,
			file_name, file_url, file_size, mime_type, created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err := r.db.ExecContext(ctx, query,
		attachment.ID,
		attachment.OrganizationID,
		attachment.UploadedBy,
		attachment.AssessmentID,
		attachment.SubmissionID,
		attachment.FileName,
		attachment.FileURL,
		attachment.FileSize,
		attachment.MIMEType,
		attachment.CreatedAt,
	)
	if err != nil {
		r.log.WithError(err).WithField("attachment_id", attachment.ID).Error("failed to insert attachment")
		return fmt.Errorf("failed to insert attachment: %w", err)
	}

	r.log.WithFields(logrus.Fields{"attachment_id": attachment.ID, "file_name": attachment.FileName}).Info("attachment created")
	return nil
}

func (r *AttachmentRepoPostgres) SoftDelete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE attachments SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL`

	result, err := r.db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		r.log.WithError(err).WithField("attachment_id", id).Error("failed to soft delete attachment")
		return fmt.Errorf("failed to soft delete attachment: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("attachment not found or already deleted")
	}

	r.log.WithField("attachment_id", id).Info("attachment soft deleted")
	return nil
}

func (r *AttachmentRepoPostgres) GetByID(ctx context.Context, id uuid.UUID) (*domain.Attachment, error) {
	query := `
		SELECT id, organization_id, uploaded_by, assessment_id, submission_id,
			   file_name, file_url, file_size, mime_type, created_at, deleted_at
		FROM attachments
		WHERE id = $1 AND deleted_at IS NULL`

	var a domain.Attachment
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&a.ID, &a.OrganizationID, &a.UploadedBy,
		&a.AssessmentID, &a.SubmissionID,
		&a.FileName, &a.FileURL, &a.FileSize, &a.MIMEType,
		&a.CreatedAt, &a.DeletedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		r.log.WithError(err).WithField("attachment_id", id).Error("failed to get attachment")
		return nil, fmt.Errorf("failed to get attachment: %w", err)
	}

	return &a, nil
}

func (r *AttachmentRepoPostgres) ListByAssessment(ctx context.Context, assessmentID uuid.UUID) ([]domain.Attachment, error) {
	return r.listBy(ctx, "assessment_id", assessmentID)
}

func (r *AttachmentRepoPostgres) ListBySubmission(ctx context.Context, submissionID uuid.UUID) ([]domain.Attachment, error) {
	return r.listBy(ctx, "submission_id", submissionID)
}

func (r *AttachmentRepoPostgres) listBy(ctx context.Context, column string, id uuid.UUID) ([]domain.Attachment, error) {
	query := fmt.Sprintf(`
		SELECT id, organization_id, uploaded_by, assessment_id, submission_id,
			   file_name, file_url, file_size, mime_type, created_at
		FROM attachments
		WHERE %s = $1 AND deleted_at IS NULL
		ORDER BY created_at ASC`, column)

	rows, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		r.log.WithError(err).WithField(column, id).Error("failed to list attachments")
		return nil, fmt.Errorf("failed to list attachments: %w", err)
	}
	defer rows.Close()

	var attachments []domain.Attachment
	for rows.Next() {
		var a domain.Attachment
		if err := rows.Scan(
			&a.ID, &a.OrganizationID, &a.UploadedBy,
			&a.AssessmentID, &a.SubmissionID,
			&a.FileName, &a.FileURL, &a.FileSize, &a.MIMEType,
			&a.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan attachment: %w", err)
		}
		attachments = append(attachments, a)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating attachments: %w", err)
	}

	return attachments, nil
}

func (r *AttachmentRepoPostgres) CountByAssessment(ctx context.Context, assessmentID uuid.UUID) (int, error) {
	return r.countBy(ctx, "assessment_id", assessmentID)
}

func (r *AttachmentRepoPostgres) CountBySubmission(ctx context.Context, submissionID uuid.UUID) (int, error) {
	return r.countBy(ctx, "submission_id", submissionID)
}

func (r *AttachmentRepoPostgres) countBy(ctx context.Context, column string, id uuid.UUID) (int, error) {
	query := fmt.Sprintf(`SELECT COUNT(*) FROM attachments WHERE %s = $1 AND deleted_at IS NULL`, column)

	var count int
	err := r.db.QueryRowContext(ctx, query, id).Scan(&count)
	if err != nil {
		r.log.WithError(err).WithField(column, id).Error("failed to count attachments")
		return 0, fmt.Errorf("failed to count attachments: %w", err)
	}

	return count, nil
}
