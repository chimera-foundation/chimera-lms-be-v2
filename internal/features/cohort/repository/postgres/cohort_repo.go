package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/cohort/domain"
	"github.com/google/uuid"
)

type CohortRepositoryPostgres struct {
	db *sql.DB
}

func NewCohortRepository(db *sql.DB) domain.CohortRepository {
	return &CohortRepositoryPostgres{
		db: db,
	}
}

func (r *CohortRepositoryPostgres) Create(ctx context.Context, cohort *domain.Cohort) error {
	query := `
		INSERT INTO cohorts (id, organization_id, academic_period_id, education_level_id, name, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	cohort.PrepareCreate(nil)

	_, err := r.db.ExecContext(ctx, query,
		cohort.ID,
		cohort.OrganizationID,
		cohort.AcademicPeriodID,
		cohort.EducationLevelID,
		cohort.Name,
		cohort.CreatedAt,
		cohort.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create cohort: %w", err)
	}

	return nil
}

func (r *CohortRepositoryPostgres) GetIDsByUserID(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	query := `
			SELECT cohort_id FROM cohort_members
			WHERE user_id = $1`

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

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *CohortRepositoryPostgres) GetByID(ctx context.Context, id uuid.UUID) (*domain.Cohort, error) {
	query := `
			SELECT id, organization_id, academic_period_id, education_level_id, name 
			FROM cohorts
			WHERE id = $1`

	cohort := &domain.Cohort{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&cohort.ID,
		&cohort.OrganizationID,
		&cohort.AcademicPeriodID,
		&cohort.EducationLevelID,
		&cohort.Name,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get cohort by id: %w", err)
	}

	return cohort, nil
}
