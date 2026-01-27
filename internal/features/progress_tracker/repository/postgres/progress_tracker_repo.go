package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/progress_tracker/domain"
)

type ProgressTrackerRepoPostgres struct {
	db *sql.DB
}

func NewProgressTrackerRepository(db *sql.DB) domain.ProgressTrackerRepository {
	return &ProgressTrackerRepoPostgres{db: db}
}

func (r *ProgressTrackerRepoPostgres) Create(ctx context.Context, tracker *domain.ProgressTracker) error {
	query := `
		INSERT INTO progress_trackers (id, enrollment_id, content_id, is_completed, updated_at)
		VALUES ($1, $2, $3, $4, $5)`

	tracker.PrepareCreate(nil)

	_, err := r.db.ExecContext(ctx, query,
		tracker.ID,
		tracker.EnrollmentID,
		tracker.ContentID,
		tracker.IsCompleted,
		tracker.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create progress tracker: %w", err)
	}

	return nil
}

func (r *ProgressTrackerRepoPostgres) Update(ctx context.Context, tracker *domain.ProgressTracker) error {
	query := `
		UPDATE progress_trackers
		SET is_completed = $2, updated_at = $3
		WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query,
		tracker.ID,
		tracker.IsCompleted,
		tracker.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update progress tracker: %w", err)
	}

	return nil
}
