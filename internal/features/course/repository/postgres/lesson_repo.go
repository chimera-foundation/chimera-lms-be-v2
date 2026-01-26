package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/course/domain"
	"github.com/google/uuid"
)

type LessonRepoPostgres struct {
	db *sql.DB
}

func NewLessonRepository(db *sql.DB) domain.LessonRepository {
	return &LessonRepoPostgres{db: db}
}

func (r *LessonRepoPostgres) Create(ctx context.Context, lesson *domain.Lesson) error {
	query := `
		INSERT INTO lessons (id, module_id, title, order_index, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)`

	lesson.PrepareCreate(nil)

	_, err := r.db.ExecContext(ctx, query,
		lesson.ID,
		lesson.ModuleID,
		lesson.Title,
		lesson.OrderIndex,
		lesson.CreatedAt,
		lesson.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create lesson: %w", err)
	}

	return nil
}

func (r *LessonRepoPostgres) GetByID(ctx context.Context, id uuid.UUID) (*domain.Lesson, error) {
	query := `
		SELECT id, module_id, title, order_index, created_at, updated_at
		FROM lessons
		WHERE id = $1 AND deleted_at IS NULL`

	lesson := &domain.Lesson{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&lesson.ID,
		&lesson.ModuleID,
		&lesson.Title,
		&lesson.OrderIndex,
		&lesson.CreatedAt,
		&lesson.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get lesson by id: %w", err)
	}

	return lesson, nil
}
