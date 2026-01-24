package domain

import (
	"time"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/shared"
	"github.com/google/uuid"
)

type EventType string

const (
	Holiday EventType = "holiday"
	Deadline EventType = "deadline"
	Session EventType = "session"
	Vanilla EventType = "vanilla"
	Meeting EventType = "meeting"
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

	Title string
	Description string
	Location string
	EventType EventType

	Color string //hex format

	StartAt *time.Time
	EndAt *time.Time
	IsAllDay bool
	RecurrenceRule *string

	Scope EventScope
	CohortID *uuid.UUID
	SectionID *uuid.UUID
	UserID *uuid.UUID

	SourceID *uuid.UUID
	SourceType *string
}