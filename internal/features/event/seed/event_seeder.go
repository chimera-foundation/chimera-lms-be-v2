package seed

import (
	"context"
	"errors"
	"fmt"
	"time"

	courseDomain "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/course/domain"
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/event/domain"
	sectionDomain "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/section/domain"
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/shared/auth"
	"github.com/google/uuid"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

type EventSeeder struct {
	r domain.EventRepository
}

func NewEventSeeder(r domain.EventRepository) *EventSeeder {
	return &EventSeeder{
		r: r,
	}
}

// SeedIndonesiaHolidays fetches holiday events from Google Calendar API.
func (s *EventSeeder) SeedIndonesiaHolidays(ctx context.Context, year int, apiKey string) ([]*domain.Event, error) {
	orgID, ok := auth.GetOrgID(ctx)
	if !ok {
		return nil, errors.New("organization ID not found in context")
	}

	if apiKey == "" {
		return nil, errors.New("GOOGLE_API_KEY is required for holiday seeding")
	}

	srv, err := calendar.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve calendar client: %w", err)
	}

	// Calendar ID for Indonesian holidays
	calendarID := "en.indonesian#holiday@group.v.calendar.google.com"

	// Define time range for the specified year
	timeMin := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC).Format(time.RFC3339)
	timeMax := time.Date(year, 12, 31, 23, 59, 59, 0, time.UTC).Format(time.RFC3339)

	fmt.Printf("DEBUG: Fetching holidays from %s to %s for calendar %s\n", timeMin, timeMax, calendarID)

	events, err := srv.Events.List(calendarID).
		ShowDeleted(false).
		SingleEvents(true).
		TimeMin(timeMin).
		TimeMax(timeMax).
		OrderBy("startTime").
		Do()
	if err != nil {
		fmt.Printf("DEBUG: Failed to list events: %v\n", err)
		return nil, fmt.Errorf("unable to retrieve holidays from google calendar: %w", err)
	}

	fmt.Printf("DEBUG: Found %d holidays\n", len(events.Items))

	var seededEvents []*domain.Event

	for _, item := range events.Items {
		var startAt, endAt time.Time
		var isAllDay bool

		if item.Start.Date != "" {
			// All-day event
			isAllDay = true
			startAt, _ = time.Parse("2006-01-02", item.Start.Date)
			// Google Calendar end date is exclusive, so it's already the start of the next day
			endAt, _ = time.Parse("2006-01-02", item.End.Date)
		} else {
			// Timed event (though holidays are usually all-day)
			startAt, _ = time.Parse(time.RFC3339, item.Start.DateTime)
			endAt, _ = time.Parse(time.RFC3339, item.End.DateTime)
		}

		// Create global holiday event
		event := domain.NewEvent(
			orgID,
			item.Summary,
			domain.Holiday,
			domain.WithColor("#EF4444"), // Red for holidays
			domain.WithTimes(startAt, endAt),
		)

		if isAllDay {
			domain.AsAllDay()(event)
		}

		if err := s.r.Create(ctx, event); err != nil {
			return nil, fmt.Errorf("failed to create holiday event %s: %w", item.Summary, err)
		}
		seededEvents = append(seededEvents, event)
	}

	return seededEvents, nil
}

// LessonWithSection associates a lesson with its section for scheduling.
type LessonWithSection struct {
	Lesson  *courseDomain.Lesson
	Section *sectionDomain.Section
}

// SeedLessonSchedules creates schedule events for each lesson.
// Each lesson gets a 45-minute schedule, spread across weekdays with day breaks.
func (s *EventSeeder) SeedLessonSchedules(
	ctx context.Context,
	lessons []*courseDomain.Lesson,
	sectionsByModule map[uuid.UUID]*sectionDomain.Section,
) ([]*domain.Event, error) {
	orgID, ok := auth.GetOrgID(ctx)
	if !ok {
		return nil, errors.New("organization ID not found in context")
	}

	var seededEvents []*domain.Event

	// Start date: Monday, January 12, 2026, at 08:00 WIB (UTC+7)
	loc := time.FixedZone("WIB", 7*60*60)
	currentDate := time.Date(2026, time.January, 12, 8, 0, 0, 0, loc)

	lessonDuration := 45 * time.Minute

	for _, lesson := range lessons {
		// Skip weekends
		for currentDate.Weekday() == time.Saturday || currentDate.Weekday() == time.Sunday {
			currentDate = currentDate.AddDate(0, 0, 1)
		}

		section, ok := sectionsByModule[lesson.ModuleID]
		if !ok {
			// If no section mapping found, skip this lesson
			continue
		}

		startTime := currentDate
		endTime := startTime.Add(lessonDuration)

		// Create schedule event with section scope
		sourceType := "lesson"
		event := &domain.Event{
			OrganizationID: orgID,
			Title:          lesson.Title,
			EventType:      domain.Schedule,
			Scope:          domain.ScopeSection,
			SectionID:      &section.ID,
			StartAt:        &startTime,
			EndAt:          &endTime,
			Color:          "#3B82F6", // Blue for schedules
			SourceID:       &lesson.ID,
			SourceType:     &sourceType,
		}

		if err := s.r.Create(ctx, event); err != nil {
			return nil, fmt.Errorf("failed to create schedule event for lesson %s: %w", lesson.Title, err)
		}
		seededEvents = append(seededEvents, event)

		// Move to next day (day break between lessons)
		currentDate = currentDate.AddDate(0, 0, 1)
	}

	return seededEvents, nil
}
