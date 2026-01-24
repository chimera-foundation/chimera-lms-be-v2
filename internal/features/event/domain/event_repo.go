package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type EventFilter struct {
    OrganizationID uuid.UUID
    
    // Scoping (The "Who")
    UserID    *uuid.UUID
    SectionID *uuid.UUID
    CohortID  *uuid.UUID
    IncludeGlobal bool

    // Categories (The "What")
    Types []EventType // e.g., []EventType{Session, Holiday}

    // Timeline (The "When")
    StartTime time.Time
    EndTime   time.Time
}

type EventRepository interface {
	Create(ctx context.Context, event *Event) error
	Update(ctx context.Context, event *Event) error
	Delete(ctx context.Context, event_id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*Event, error)
    Find(ctx context.Context, filter EventFilter) ([]*Event, error)
}