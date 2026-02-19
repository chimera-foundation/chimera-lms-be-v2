package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

const (
	MaxFileSize             = 25 << 20 // 25 MB
	MaxAttachmentsPerEntity = 10
)

var AllowedMIMETypes = map[string]bool{
	// Documents
	"application/pdf":    true,
	"application/msword": true,
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
	"application/vnd.ms-excel": true,
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":         true,
	"application/vnd.ms-powerpoint":                                             true,
	"application/vnd.openxmlformats-officedocument.presentationml.presentation": true,
	"text/plain": true,

	// Images
	"image/png":  true,
	"image/jpeg": true,
	"image/gif":  true,

	// Archives
	"application/zip": true,
}

type Attachment struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	UploadedBy     uuid.UUID
	AssessmentID   *uuid.UUID
	SubmissionID   *uuid.UUID
	FileName       string
	FileURL        string
	FileSize       int64
	MIMEType       string
	CreatedAt      time.Time
	DeletedAt      *time.Time
}

func (a *Attachment) Validate() error {
	if a.FileName == "" {
		return errors.New("file name is required")
	}
	if a.FileSize <= 0 {
		return errors.New("file size must be positive")
	}
	if a.FileSize > MaxFileSize {
		return errors.New("file size exceeds the maximum allowed (25 MB)")
	}
	if !AllowedMIMETypes[a.MIMEType] {
		return errors.New("file type is not allowed")
	}
	if a.AssessmentID == nil && a.SubmissionID == nil {
		return errors.New("attachment must be linked to an assessment or a submission")
	}
	if a.AssessmentID != nil && a.SubmissionID != nil {
		return errors.New("attachment cannot be linked to both an assessment and a submission")
	}
	return nil
}
