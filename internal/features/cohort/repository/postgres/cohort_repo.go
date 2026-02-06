package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/cohort/domain"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type CohortRepositoryPostgres struct {
	db  *sql.DB
	log *logrus.Logger
}

func NewCohortRepository(db *sql.DB, log *logrus.Logger) domain.CohortRepository {
	return &CohortRepositoryPostgres{
		db:  db,
		log: log,
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
		r.log.WithError(err).WithField("cohort_id", cohort.ID).Error("failed to create cohort")
		return fmt.Errorf("failed to create cohort: %w", err)
	}

	r.log.WithFields(logrus.Fields{"cohort_id": cohort.ID, "name": cohort.Name}).Info("cohort created successfully")
	return nil
}

func (r *CohortRepositoryPostgres) GetIDsByUserID(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	query := `
			SELECT cohort_id FROM cohort_members
			WHERE user_id = $1`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		r.log.WithError(err).WithField("user_id", userID).Error("failed to query cohort ids by user")
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
		r.log.WithField("cohort_id", id).Debug("cohort not found by id")
		return nil, nil
	}
	if err != nil {
		r.log.WithError(err).WithField("cohort_id", id).Error("failed to get cohort by id")
		return nil, fmt.Errorf("failed to get cohort by id: %w", err)
	}

	return cohort, nil
}
