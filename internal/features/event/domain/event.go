package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/shared"
	"github.com/google/uuid"
)

type EventType string

const (
	Holiday  EventType = "holiday"
	Deadline EventType = "deadline"
	Session  EventType = "session"
	Vanilla  EventType = "vanilla"
	Meeting  EventType = "meeting"
	Schedule EventType = "schedule"
)

type EventScope string

const (
	ScopeGlobal   EventScope = "global"
	ScopeCohort   EventScope = "cohort"
	ScopeSection  EventScope = "section"
	ScopePersonal EventScope = "personal"
)

type Event struct {
	shared.Base

	OrganizationID uuid.UUID

	Title       string
	Description string
	Location    string
	EventType   EventType

	Color string //hex format

	StartAt        *time.Time
	EndAt          *time.Time
	IsAllDay       bool
	RecurrenceRule *string

	Scope     EventScope
	CohortID  *uuid.UUID
	SectionID *uuid.UUID
	UserID    *uuid.UUID

	SourceID   *uuid.UUID
	SourceType *string

	ImageURL *string
}

type EventOption func(*Event)

func Recurring(rrule string) EventOption {
	return func(e *Event) { e.RecurrenceRule = &rrule }
}

func AsAllDay() EventOption {
	return func(e *Event) { e.IsAllDay = true }
}

func ForCohort(id uuid.UUID) EventOption {
	return func(e *Event) {
		e.Scope = ScopeCohort
		e.CohortID = &id
	}
}

func ForSection(id uuid.UUID) EventOption {
	return func(e *Event) {
		e.Scope = ScopeSection
		e.SectionID = &id
	}
}

func ForUser(id uuid.UUID) EventOption {
	return func(e *Event) {
		e.Scope = ScopePersonal
		e.UserID = &id
	}
}

func WithTimes(start, end time.Time) EventOption {
	return func(e *Event) {
		e.StartAt = &start
		e.EndAt = &end
	}
}

func WithImage(url string) EventOption {
	return func(e *Event) { e.ImageURL = &url }
}

func LinkedTo(id uuid.UUID, entityType string) EventOption {
	return func(e *Event) {
		e.SourceID = &id
		e.SourceType = &entityType
	}
}

func WithColor(color string) EventOption {
	return func(e *Event) { e.Color = color }
}

func WithLocation(loc string) EventOption {
	return func(e *Event) { e.Location = loc }
}

func NewEvent(orgID uuid.UUID, title string, eventType EventType, opts ...EventOption) *Event {
	e := &Event{
		OrganizationID: orgID,
		Title:          title,
		EventType:      eventType,
		Scope:          ScopeGlobal,
		Color:          "#3B82F6",
	}

	for _, opt := range opts {
		opt(e)
	}
	return e
}

func (e *Event) Validate() error {
	if e.OrganizationID == uuid.Nil {
		return errors.New("organization_id is required")
	}
	if strings.TrimSpace(e.Title) == "" {
		return errors.New("event title cannot be empty")
	}

	switch e.EventType {
	case Session, Meeting, Schedule:
		if e.StartAt == nil || e.EndAt == nil {
			return errors.New("events must have both a start and end time")
		}
		if e.EndAt.Before(*e.StartAt) {
			return errors.New("end time cannot be before start time")
		}

	case Deadline:
		if e.StartAt == nil {
			return errors.New("deadline events must have a start_at (due date)")
		}
	}

	switch e.Scope {
	case ScopeSection:
		if e.SectionID == nil || *e.SectionID == uuid.Nil {
			return errors.New("section scope requires a valid section_id")
		}
	case ScopeCohort:
		if e.CohortID == nil || *e.CohortID == uuid.Nil {
			return errors.New("cohort scope requires a valid cohort_id")
		}
	case ScopePersonal:
		if e.UserID == nil || *e.UserID == uuid.Nil {
			return errors.New("personal scope requires a valid user_id")
		}
	}

	if e.SourceID != nil && (e.SourceType == nil || *e.SourceType == "") {
		return errors.New("source_id is provided but source_type is missing")
	}

	return nil
}
