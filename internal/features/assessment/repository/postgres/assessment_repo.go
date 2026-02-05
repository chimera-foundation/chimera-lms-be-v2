package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/assessment/domain"
)

type AssessmentRepoPostgres struct {
	db *sql.DB
}

func NewAssessmentRepoPostgres(db *sql.DB) domain.AssessmentRepo {
	return &AssessmentRepoPostgres{
		db: db,
	}
}

func (r *AssessmentRepoPostgres) Create(ctx context.Context, assessment *domain.Assessment) error {
	query := `
			INSERT INTO assessments (
				id,
				organization_id,
				title,
				assessment_type,
				assessment_sub_type,
				due_date,
				created_at, 
				updated_at, 
				course_id
			) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := r.db.ExecContext(ctx, query, 
		assessment.ID,
		assessment.OrganizationID,
		assessment.Title,
		assessment.Type,
		assessment.SubType,
		assessment.DueDate,
		assessment.CreatedAt,
		assessment.UpdatedAt,
		assessment.CourseID,
	)

	if err != nil {
		return fmt.Errorf("failed to insert assessment: %w", err)
	}

	return nil
}