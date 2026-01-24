package service

import (
	"context"
	"errors"
	"time"

	c "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/cohort/domain"
	en "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/enrollment/domain"
	e "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/event/domain"
	o "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/organization/domain"
	s "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/section/domain"
	"github.com/google/uuid"
)

type eventService struct {
	repo           e.EventRepository
	orgRepo        o.OrganizationRepository
	enrollmentRepo en.EnrollmentRepository
	cohortRepo     c.CohortRepository
	sectionRepo    s.SectionRepository // for staff/teachers
}

func NewEventService(
	repo           e.EventRepository,
	orgRepo        o.OrganizationRepository,
	enrollmentRepo en.EnrollmentRepository,
	cohortRepo     c.CohortRepository,
	sectionRepo    s.SectionRepository,
) EventService {
	return &eventService{
		repo:           repo,
		orgRepo:        orgRepo,
		enrollmentRepo: enrollmentRepo,
		cohortRepo:     cohortRepo,
		sectionRepo:    sectionRepo,
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

	return s.repo.Find(ctx, filter)
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

    if err := s.repo.Update(ctx, e); err != nil {
        return nil, err
    }

    return e, nil
}

func (s *eventService) DeleteEvent(ctx context.Context, id uuid.UUID) error {
    // TODO: Check permissions here
    return s.repo.Delete(ctx, id)
}