package shared

import (
	"time"
	"github.com/google/uuid"
)

type Base struct {
	ID uuid.UUID `validate:"required"`
	
	// Audit metada
	CreatedAt time.Time
	CreatedBy uuid.UUID
	UpdatedAt time.Time
	UpdatedBy uuid.UUID

	// Optional/nullable fields
	DeletedAt *time.Time
	DeletedBy *uuid.UUID
}