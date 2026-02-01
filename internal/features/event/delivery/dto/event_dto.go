package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateEventRequest struct {
	Title          string     `json:"title"`
	Description    string     `json:"description"`
	Location       string     `json:"location"`
	EventType      string     `json:"event_type"` // Maps to e.EventType
	Color          string     `json:"color"`
	StartAt        *time.Time `json:"start_at"`
	EndAt          *time.Time `json:"end_at"`
	IsAllDay       bool       `json:"is_all_day"`
	RecurrenceRule *string    `json:"recurrence_rule"`
	Scope          string     `json:"scope"` // Maps to e.EventScope
	CohortID       *uuid.UUID `json:"cohort_id"`
	SectionID      *uuid.UUID `json:"section_id"`
	UserID         *uuid.UUID `json:"user_id"`
	ImageURL       *string    `json:"image_url"`
}