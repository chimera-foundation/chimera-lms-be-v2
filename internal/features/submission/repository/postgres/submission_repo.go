package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/submission/domain"
)

type SubmissionRepoPostgres struct {
	db *sql.DB
}

func NewSubmissionRepoPostgres(db *sql.DB) domain.SubmissionRepository {
	return &SubmissionRepoPostgres{
		db: db,
	}
}

func (r *SubmissionRepoPostgres) Create(ctx context.Context, submission *domain.Submission) error {
	query := `
		INSERT INTO submissions (
			id,
			assessment_id,
			user_id,
			enrollment_id,
			final_score,
			submitted_at,
			created_at,
			updated_at
		) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := r.db.ExecContext(ctx, query,
		submission.ID,
		submission.AssessmentID,
		submission.UserID,
		submission.EnrollmentID,
		submission.FinalScore,
		submission.SubmittedAt,
		submission.CreatedAt,
		submission.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to insert submission: %w", err)
	}

	return nil
}
