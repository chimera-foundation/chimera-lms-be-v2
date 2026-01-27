package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/enrollment/domain"
	"github.com/google/uuid"
)

type EnrollmentRepositoryPostgres struct {
	db *sql.DB
}

func NewEnrollmentRepository(db *sql.DB) domain.EnrollmentRepository {
	return &EnrollmentRepositoryPostgres{
		db: db,
	}
}

func (r *EnrollmentRepositoryPostgres) Create(ctx context.Context, enrollment *domain.Enrollment) error {
	query := `
		INSERT INTO enrollments (id, user_id, course_id, section_id, academic_period_id, status, enrolled_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	enrollment.PrepareCreate(nil)

	_, err := r.db.ExecContext(ctx, query,
		enrollment.ID,
		enrollment.UserID,
		enrollment.CourseID,
		enrollment.SectionID,
		enrollment.AcademicPeriodID,
		enrollment.Status,
		enrollment.EnrolledAt,
		enrollment.CreatedAt,
		enrollment.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create enrollment: %w", err)
	}

	return nil
}

func (r *EnrollmentRepositoryPostgres) GetActiveSectionIDsByUserID(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	query := `
			SELECT section_id 
			FROM enrollments 
			WHERE user_id = $1 AND status = 'active'`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []uuid.UUID
	for rows.Next() {
		var id uuid.UUID

		if err := rows.Scan(&id); err != nil {
			return nil, err
		}

		result = append(result, id)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
