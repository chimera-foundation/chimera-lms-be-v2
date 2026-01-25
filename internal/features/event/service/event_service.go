package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	c "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/cohort/domain"
	en "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/enrollment/domain"
	e "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/event/domain"
	o "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/organization/domain"
	s "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/section/domain"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type eventService struct {
	repo           e.EventRepository
	orgRepo        o.OrganizationRepository
	enrollmentRepo en.EnrollmentRepository
	cohortRepo     c.CohortRepository
	sectionRepo    s.SectionRepository // for staff/teachers
	redis *redis.Client
}

func NewEventService(
	repo           e.EventRepository,
	orgRepo        o.OrganizationRepository,
	enrollmentRepo en.EnrollmentRepository,
	cohortRepo     c.CohortRepository,
	sectionRepo    s.SectionRepository,
	redis *redis.Client,
) EventService {
	return &eventService{
		repo:           repo,
		orgRepo:        orgRepo,
		enrollmentRepo: enrollmentRepo,
		cohortRepo:     cohortRepo,
		sectionRepo:    sectionRepo,
		redis: redis,
	}
}

func (s *eventService) CreateEvent(ctx context.Context, e *e.Event) (*e.Event, error) {
	if err := e.Validate(); err != nil {
		return nil, err
	}

	// 2. Security/Context check: Ensure the Org exists
	// (TODO: verify the creator has permission for this Org)
	
	err := s.repo.Create(ctx, e)
	if err != nil {
		return nil, err
	}

	return e, nil
}

func (s *eventService) GetCalendarForUser(ctx context.Context, userID uuid.UUID, start, end time.Time) ([]*e.Event, error) {
	cacheKey := fmt.Sprintf("events:cal:%s:%d:%d", userID, start.Unix(), end.Unix())

	val, err := s.redis.Get(ctx, cacheKey).Result()
    if err == nil {
        var cachedEvents []*e.Event
        if err := json.Unmarshal([]byte(val), &cachedEvents); err == nil {
            return cachedEvents, nil
        }
    }

	orgID, err := s.orgRepo.GetIDByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	
	activeEnrollments, _ := s.enrollmentRepo.GetActiveSectionIDsByUserID(ctx, userID)
	staffSections, _ := s.sectionRepo.GetSectionIDsByUserID(ctx, userID)
	
	// Merge unique IDs
	sectionIDMap := make(map[uuid.UUID]bool)
	for _, id := range activeEnrollments { sectionIDMap[id] = true }
	for _, id := range staffSections { sectionIDMap[id] = true }
	
	var sectionIDs []uuid.UUID
	for id := range sectionIDMap {
		sectionIDs = append(sectionIDs, id)
	}

	// 3. Identify Cohorts
	cohortIDs, _ := s.cohortRepo.GetIDsByUserID(ctx, userID)

	// 4. Execute Scoped Search
	filter := e.EventFilter{
		OrganizationID: orgID,
		UserID:         &userID,
		SectionIDs:     sectionIDs,
		CohortIDs:      cohortIDs,
		IncludeGlobal:  true,
		StartTime:      start,
		EndTime:        end,
	}

	events, err := s.repo.Find(ctx, filter)
    if err != nil {
        return nil, err
    }

    go func(evs []*e.Event) {
        data, _ := json.Marshal(evs)
        s.redis.Set(context.Background(), cacheKey, data, 15*time.Minute)
    }(events)

	return events, nil
}

func (s *eventService) GetSectionSchedule(ctx context.Context, sectionID uuid.UUID, start, end time.Time) ([]*e.Event, error) {
    // 1. Get the Section to find its parent Cohort
    section, err := s.sectionRepo.GetByID(ctx, sectionID)
    if err != nil {
        return nil, err
    }

    // 2. Get the Cohort to find the OrganizationID
    // Your schema confirms cohorts table has organization_id
    cohort, err := s.cohortRepo.GetByID(ctx, section.CohortID)
    if err != nil {
        return nil, err
    }

    // 3. Now we have the OrganizationID required by the EventFilter
    filter := e.EventFilter{
        OrganizationID: cohort.OrganizationID,
        SectionIDs:     []uuid.UUID{sectionID},
        IncludeGlobal:  false,
        StartTime:      start,
        EndTime:        end,
        Limit:          100,
    }

    return s.repo.Find(ctx, filter)
}

func (s *eventService) GetAnnouncements(ctx context.Context, orgID uuid.UUID, limit, offset int) ([]*e.Event, error) {
    if limit <= 0 { limit = 10 }

    filter := e.EventFilter{
        OrganizationID: orgID,
        Types:          []e.EventType{e.Announcement},
        IncludeGlobal:  true, // Announcements are usually global or cohort-wide
        Limit:          limit,
        Offset:         offset,
    }

    return s.repo.Find(ctx, filter)
}

func (s *eventService) GetEvents(ctx context.Context, orgID uuid.UUID, limit, offset int) ([]*e.Event, error) {
    if limit <= 0 { limit = 20 }

    filter := e.EventFilter{
        OrganizationID: orgID,
        IncludeGlobal:  true,
        Limit:          limit,
        Offset:         offset,
    }

    return s.repo.Find(ctx, filter)
}

func (s *eventService) flushUserCache(ctx context.Context, userID uuid.UUID) {
	pattern := fmt.Sprintf("events:cal:%s:*", userID)
	iter := s.redis.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		s.redis.Del(ctx, iter.Val())
	}
}

func (s *eventService) UpdateEvent(ctx context.Context, e *e.Event) (*e.Event, error) {
    if err := e.Validate(); err != nil {
        return nil, err
    }

    // Ensure the event exists before updating
    existing, err := s.repo.GetByID(ctx, e.ID)
    if err != nil || existing == nil {
        return nil, errors.New("event not found")
    }

    // Business Logic: Prevent moving events across Organizations
    if existing.OrganizationID != e.OrganizationID {
        return nil, errors.New("unauthorized: organization mismatch")
    }

	err = s.repo.Update(ctx, e) 
	if err == nil && e.UserID != nil {
        s.flushUserCache(ctx, *e.UserID)
    }

    return e, nil
}

func (s *eventService) DeleteEvent(ctx context.Context, id uuid.UUID) error {
    // TODO: Check permissions here
    return s.repo.Delete(ctx, id)
}