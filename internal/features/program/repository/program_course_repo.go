package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/program/domain"
)

type ProgramCourseRepoPostgres struct {
	db *sql.DB
}

func NewProgramCourseRepository(db *sql.DB) domain.ProgramCourseRepository {
	return &ProgramCourseRepoPostgres{db: db}
}

func (r *ProgramCourseRepoPostgres) Create(ctx context.Context, pc *domain.ProgramCourse) error {
	query := `
		INSERT INTO program_courses (program_id, course_id, order_index)
		VALUES ($1, $2, $3)
		ON CONFLICT DO NOTHING`

	_, err := r.db.ExecContext(ctx, query,
		pc.ProgramID,
		pc.CourseID,
		pc.OrderIndex,
	)

	if err != nil {
		return fmt.Errorf("failed to create program course: %w", err)
	}

	return nil
}
