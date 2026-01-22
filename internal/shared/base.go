package shared

import (
	"time"
	"github.com/google/uuid"
)

type Base struct {
	ID uuid.UUID `validate:"required"`
	
	// Audit metada
	CreatedAt time.Time
	CreatedBy *uuid.UUID
	UpdatedAt time.Time
	UpdatedBy *uuid.UUID

	// Optional/nullable fields
	DeletedAt *time.Time
	DeletedBy *uuid.UUID
}

func (b *Base) PrepareCreate(actorID *uuid.UUID) {
    if b.ID == uuid.Nil {
        b.ID = uuid.New()
    }
    now := time.Now()
    b.CreatedAt = now
    b.UpdatedAt = now
    b.CreatedBy = actorID
    b.UpdatedBy = actorID
}