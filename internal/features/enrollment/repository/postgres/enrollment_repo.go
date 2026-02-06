package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/enrollment/domain"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type EnrollmentRepositoryPostgres struct {
	db  *sql.DB
	log *logrus.Logger
}

func NewEnrollmentRepository(db *sql.DB, log *logrus.Logger) domain.EnrollmentRepository {
	return &EnrollmentRepositoryPostgres{
		db:  db,
		log: log,
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
		r.log.WithError(err).WithField("enrollment_id", enrollment.ID).Error("failed to create enrollment")
		return fmt.Errorf("failed to create enrollment: %w", err)
	}

	r.log.WithFields(logrus.Fields{"enrollment_id": enrollment.ID, "user_id": enrollment.UserID}).Info("enrollment created successfully")
	return nil
}

func (r *EnrollmentRepositoryPostgres) GetActiveSectionIDsByUserID(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	query := `
			SELECT section_id 
			FROM enrollments 
			WHERE user_id = $1 AND status = 'active'`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		r.log.WithError(err).WithField("user_id", userID).Error("failed to query active section ids")
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
		r.log.WithError(err).WithField("user_id", userID).Error("error iterating active section ids")
		return nil, err
	}

	return result, nil
}
