package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type EventFilter struct {
    OrganizationID uuid.UUID
    
    Limit  int // How many to fetch
    Offset int // How many to skip
    
    // Scoping (The "Who")
    UserID    *uuid.UUID
    SectionIDs []uuid.UUID  // Changed from *uuid.UUID (Classroom scope)
    CohortIDs  []uuid.UUID  // Changed from *uuid.UUID (Year Group scope)
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
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*Event, error)
    Find(ctx context.Context, filter EventFilter) ([]*Event, error)
}