package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/course/domain"
	"github.com/google/uuid"
)

type CourseRepoPostgres struct {
	db *sql.DB
}

func NewCourseRepository(db *sql.DB) domain.CourseRepository {
	return &CourseRepoPostgres{db: db}
}

func (r *CourseRepoPostgres) Create(ctx context.Context, course *domain.Course) error {
	query := `
		INSERT INTO courses (id, organization_id, instructor_id, subject_id, education_level_id, title, description, status, price, grade_level, credits, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`

	course.PrepareCreate(nil)

	_, err := r.db.ExecContext(ctx, query,
		course.ID,
		course.OrganizationID,
		course.InstructorID,
		course.SubjectID,
		course.EducationLevelID,
		course.Title,
		course.Description,
		course.Status,
		course.Price,
		course.GradeLevel,
		course.Credits,
		course.CreatedAt,
		course.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create course: %w", err)
	}

	return nil
}

func (r *CourseRepoPostgres) GetByID(ctx context.Context, id uuid.UUID) (*domain.Course, error) {
	query := `
		SELECT id, organization_id, instructor_id, subject_id, education_level_id, title, description, status, price, grade_level, credits, created_at, updated_at
		FROM courses
		WHERE id = $1 AND deleted_at IS NULL`

	course := &domain.Course{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&course.ID,
		&course.OrganizationID,
		&course.InstructorID,
		&course.SubjectID,
		&course.EducationLevelID,
		&course.Title,
		&course.Description,
		&course.Status,
		&course.Price,
		&course.GradeLevel,
		&course.Credits,
		&course.CreatedAt,
		&course.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get course by id: %w", err)
	}

	return course, nil
}
