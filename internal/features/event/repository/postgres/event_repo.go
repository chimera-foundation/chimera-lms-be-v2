package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/event/domain"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type EventRepoPostgres struct {
	db  *sql.DB
	log *logrus.Logger
}

func NewEventRepository(db *sql.DB, log *logrus.Logger) domain.EventRepository {
	return &EventRepoPostgres{
		db:  db,
		log: log,
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

	if err != nil {
		r.log.WithError(err).WithField("event_id", e.ID).Error("failed to create event")
		return err
	}

	r.log.WithFields(logrus.Fields{"event_id": e.ID, "title": e.Title}).Info("event created successfully")
	return nil
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
		r.log.WithError(err).WithField("event_id", e.ID).Error("failed to update event")
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		r.log.WithField("event_id", e.ID).Warn("event not found for update")
		return errors.New("event not found")
	}

	r.log.WithField("event_id", e.ID).Info("event updated successfully")
	return nil
}

func (r *EventRepoPostgres) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE events SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		r.log.WithError(err).WithField("event_id", id).Error("failed to delete event")
		return err
	}
	r.log.WithField("event_id", id).Info("event deleted successfully")
	return nil
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
	if f.Limit <= 0 {
		f.Limit = 50
	}

	baseFields := `id, organization_id, title, description, location, event_type,
                   color, start_at, end_at, is_all_day, recurrence_rule,
                   scope, cohort_id, section_id, user_id,
                   source_id, source_type, image_url, created_at, updated_at`

	commonFilter := ` AND organization_id = $1 AND deleted_at IS NULL`

	timeAndTypeFilter := ` AND (start_at <= $2 AND (end_at >= $3 OR end_at IS NULL))`
	args := []any{f.OrganizationID, f.EndTime, f.StartTime}
	placeholderID := 4

	if len(f.Types) > 0 {
		timeAndTypeFilter += fmt.Sprintf(" AND event_type = ANY($%d)", placeholderID)
		args = append(args, pq.Array(f.Types))
		placeholderID++
	}

	var subQueries []string

	if f.IncludeGlobal {
		subQueries = append(subQueries, fmt.Sprintf("SELECT %s FROM events WHERE scope = 'global' %s %s",
			baseFields, commonFilter, timeAndTypeFilter))
	}

	if f.UserID != nil && *f.UserID != uuid.Nil {
		subQueries = append(subQueries, fmt.Sprintf("SELECT %s FROM events WHERE user_id = $%d %s %s",
			baseFields, placeholderID, commonFilter, timeAndTypeFilter))
		args = append(args, f.UserID)
		placeholderID++
	}

	if len(f.SectionIDs) > 0 {
		subQueries = append(subQueries, fmt.Sprintf("SELECT %s FROM events WHERE section_id = ANY($%d) %s %s",
			baseFields, placeholderID, commonFilter, timeAndTypeFilter))
		args = append(args, pq.Array(f.SectionIDs))
		placeholderID++
	}

	if len(f.CohortIDs) > 0 {
		subQueries = append(subQueries, fmt.Sprintf("SELECT %s FROM events WHERE cohort_id = ANY($%d) %s %s",
			baseFields, placeholderID, commonFilter, timeAndTypeFilter))
		args = append(args, pq.Array(f.CohortIDs))
		placeholderID++
	}

	if len(subQueries) == 0 {
		return []*domain.Event{}, nil
	}

	finalQuery := fmt.Sprintf("(%s) ORDER BY start_at ASC, title ASC LIMIT $%d OFFSET $%d",
		strings.Join(subQueries, ") UNION ALL ("), placeholderID, placeholderID+1)
	args = append(args, f.Limit, f.Offset)

	rows, err := r.db.QueryContext(ctx, finalQuery, args...)
	if err != nil {
		r.log.WithError(err).WithField("org_id", f.OrganizationID).Error("union find events failed")
		return nil, fmt.Errorf("union find events failed: %w", err)
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
