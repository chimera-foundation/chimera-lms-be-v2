package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/content/domain"
	"github.com/google/uuid"
)

type ContentRepoPostgres struct {
	db *sql.DB
}

func NewContentRepository(db *sql.DB) domain.ContentRepository {
	return &ContentRepoPostgres{db: db}
}

func (r *ContentRepoPostgres) Create(ctx context.Context, content *domain.Content) error {
	query := `
		INSERT INTO contents (id, lesson_id, assessment_id, content_type, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)`

	content.PrepareCreate(nil)

	var assessmentID interface{}
	if content.AssessmentID != uuid.Nil {
		assessmentID = content.AssessmentID
	}

	var lessonID interface{}
	if content.LessonID != uuid.Nil {
		lessonID = content.LessonID
	}

	_, err := r.db.ExecContext(ctx, query,
		content.ID,
		lessonID,
		assessmentID,
		content.Type,
		content.CreatedAt,
		content.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create content: %w", err)
	}

	return nil
}

func (r *ContentRepoPostgres) GetByLessonID(ctx context.Context, lessonID uuid.UUID) ([]*domain.Content, error) {
	query := `
		SELECT id, lesson_id, assessment_id, content_type, created_at, updated_at
		FROM contents
		WHERE lesson_id = $1 AND deleted_at IS NULL`

	rows, err := r.db.QueryContext(ctx, query, lessonID)
	if err != nil {
		return nil, fmt.Errorf("failed to get contents by lesson id: %w", err)
	}
	defer rows.Close()

	var contents []*domain.Content
	for rows.Next() {
		content := &domain.Content{}
		var assessmentID, lessonIDVal sql.NullString
		err := rows.Scan(
			&content.ID,
			&lessonIDVal,
			&assessmentID,
			&content.Type,
			&content.CreatedAt,
			&content.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan content: %w", err)
		}
		if lessonIDVal.Valid {
			content.LessonID, _ = uuid.Parse(lessonIDVal.String)
		}
		if assessmentID.Valid {
			content.AssessmentID, _ = uuid.Parse(assessmentID.String)
		}
		contents = append(contents, content)
	}

	return contents, nil
}
