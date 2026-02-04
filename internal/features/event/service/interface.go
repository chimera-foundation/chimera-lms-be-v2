package service

import (
	"context"
	"time"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/event/domain"
	"github.com/google/uuid"
)

type EventService interface {
	CreateEvent(ctx context.Context, e *domain.Event) (*domain.Event, error)
	UpdateEvent(ctx context.Context, e *domain.Event) (*domain.Event, error)
	DeleteEvent(ctx context.Context, id uuid.UUID) error
	GetCalendarForUser(ctx context.Context, userID uuid.UUID, start, end time.Time) ([]*domain.Event, error)
	GetSectionSchedule(ctx context.Context, sectionID uuid.UUID, start, end time.Time) ([]*domain.Event, error)
	GetAnnouncements(ctx context.Context, orgID uuid.UUID, start, end time.Time, limit int, offset int) ([]*domain.Event, error)
	GetEvents(ctx context.Context, orgID uuid.UUID, start, end time.Time, limit int, offset int) ([]*domain.Event, error)
}
