package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/organization/domain"
	"github.com/google/uuid"
)

type AcademicPeriodRepoPostgres struct {
	db *sql.DB
}

func NewAcademicPeriodRepository(db *sql.DB) domain.AcademicPeriodRepository {
	return &AcademicPeriodRepoPostgres{db: db}
}

func (r *AcademicPeriodRepoPostgres) Create(ctx context.Context, period *domain.AcademicPeriod, orgID uuid.UUID) error {
	query := `
		INSERT INTO academic_periods (id, organization_id, name, start_date, end_date, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	period.PrepareCreate(nil)

	_, err := r.db.ExecContext(ctx, query,
		period.ID,
		orgID,
		period.Name,
		period.StartDate,
		period.EndDate,
		period.IsActive,
		period.CreatedAt,
		period.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create academic period: %w", err)
	}

	return nil
}

func (r *AcademicPeriodRepoPostgres) GetActiveByOrganizationID(ctx context.Context, orgID uuid.UUID) (*domain.AcademicPeriod, error) {
	query := `
		SELECT id, name, start_date, end_date, is_active, created_at, updated_at
		FROM academic_periods
		WHERE organization_id = $1 AND is_active = true AND deleted_at IS NULL
		LIMIT 1`

	period := &domain.AcademicPeriod{}
	err := r.db.QueryRowContext(ctx, query, orgID).Scan(
		&period.ID,
		&period.Name,
		&period.StartDate,
		&period.EndDate,
		&period.IsActive,
		&period.CreatedAt,
		&period.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get active academic period: %w", err)
	}

	return period, nil
}
