package dto

import (
	"time"

	"github.com/google/uuid"
)

type AttachmentResponse struct {
	ID        uuid.UUID `json:"id"`
	FileName  string    `json:"file_name"`
	FileURL   string    `json:"file_url"`
	FileSize  int64     `json:"file_size"`
	MIMEType  string    `json:"mime_type"`
	CreatedAt time.Time `json:"created_at"`
}
