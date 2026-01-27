package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/event/domain"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type EventRepoPostgres struct {
	db *sql.DB
}

func NewEventRepository(db *sql.DB) domain.EventRepository {
	return &EventRepoPostgres{
		db: db,
	}
}

func (r *EventRepoPostgres) Create(ctx context.Context, e *domain.Event) error {
	query := `
		INSERT INTO events (
			id, organization_id, title, description, location, event_type,
			color, start_at, end_at, is_all_day, recurrence_rule,
			scope, cohort_id, section_id, user_id,
			source_id, source_type, image_url, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)`

	e.PrepareCreate(&e.OrganizationID)

	_, err := r.db.ExecContext(ctx, query,
		e.ID, e.OrganizationID, e.Title, e.Description, e.Location, e.EventType,
		e.Color, e.StartAt, e.EndAt, e.IsAllDay, e.RecurrenceRule,
		e.Scope, e.CohortID, e.SectionID, e.UserID,
		e.SourceID, e.SourceType, e.ImageURL, e.CreatedAt, e.UpdatedAt,
	)

	return err
}

func (r *EventRepoPostgres) Update(ctx context.Context, e *domain.Event) error {
	query := `
		UPDATE events SET
			title = $1, description = $2, location = $3, event_type = $4,
			color = $5, start_at = $6, end_at = $7, is_all_day = $8,
			recurrence_rule = $9, scope = $10, cohort_id = $11,
			section_id = $12, user_id = $13, source_id = $14,
			source_type = $15, image_url = $16, updated_at = $17
		WHERE id = $18 AND deleted_at IS NULL`

	res, err := r.db.ExecContext(ctx, query,
		e.Title, e.Description, e.Location, e.EventType,
		e.Color, e.StartAt, e.EndAt, e.IsAllDay,
		e.RecurrenceRule, e.Scope, e.CohortID,
		e.SectionID, e.UserID, e.SourceID,
		e.SourceType, e.ImageURL, time.Now(), e.ID,
	)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New("event not found")
	}
	return nil
}

func (r *EventRepoPostgres) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE events SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, query, time.Now(), id)
	return err
}

func (r *EventRepoPostgres) scanRow(scanner interface{ Scan(dest ...any) error }) (*domain.Event, error) {
	e := &domain.Event{}
	err := scanner.Scan(
		&e.ID, &e.OrganizationID, &e.Title, &e.Description, &e.Location, &e.EventType,
		&e.Color, &e.StartAt, &e.EndAt, &e.IsAllDay, &e.RecurrenceRule,
		&e.Scope, &e.CohortID, &e.SectionID, &e.UserID,
		&e.SourceID, &e.SourceType, &e.ImageURL, &e.CreatedAt, &e.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (r *EventRepoPostgres) GetByID(ctx context.Context, id uuid.UUID) (*domain.Event, error) {
	query := `SELECT id, organization_id, title, description, location, event_type,
					color, start_at, end_at, is_all_day, recurrence_rule,
					scope, cohort_id, section_id, user_id,
					source_id, source_type, image_url, created_at, updated_at 
			  FROM events WHERE id = $1 AND deleted_at IS NULL`
	return r.scanRow(r.db.QueryRowContext(ctx, query, id))
}

func (r *EventRepoPostgres) Find(ctx context.Context, f domain.EventFilter) ([]*domain.Event, error) {
	query := `
		SELECT id, organization_id, title, description, location, event_type,
		       color, start_at, end_at, is_all_day, recurrence_rule,
		       scope, cohort_id, section_id, user_id,
		       source_id, source_type, image_url, created_at, updated_at
		FROM events
		WHERE organization_id = $1
		  AND deleted_at IS NULL
		  AND (
		      (scope = 'global' AND $2 = true) OR
		      (scope = 'personal' AND user_id = $3) OR
		      (scope = 'section' AND section_id = ANY($4)) OR
		      (scope = 'cohort' AND cohort_id = ANY($5))
		  )
		  AND (event_type = ANY($6) OR $7 = 0)
		  AND (start_at <= $8 AND (end_at >= $9 OR end_at IS NULL))
		ORDER BY start_at ASC, title ASC 
		LIMIT $10 OFFSET $11`

	// Handle the any-type logic
	typeCount := len(f.Types)

	if f.Limit == 0 {
		f.Limit = 50 
	}

	rows, err := r.db.QueryContext(ctx, query,
		f.OrganizationID, f.IncludeGlobal, f.UserID, pq.Array(f.SectionIDs), pq.Array(f.CohortIDs),
		pq.Array(f.Types), typeCount, f.EndTime, f.StartTime, f.Limit, f.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*domain.Event
	for rows.Next() {
		e, err := r.scanRow(rows)
		if err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, nil
}
