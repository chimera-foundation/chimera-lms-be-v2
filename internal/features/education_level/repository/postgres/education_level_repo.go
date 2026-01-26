package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/education_level/domain"
	"github.com/google/uuid"
)

type EducationLevelRepoPostgres struct {
	db *sql.DB
}

func NewEducationLevelRepository(db *sql.DB) domain.EducationLevelRepository {
	return &EducationLevelRepoPostgres{db: db}
}

func (r *EducationLevelRepoPostgres) Create(ctx context.Context, level *domain.EducationLevel) error {
	query := `
		INSERT INTO education_levels (id, organization_id, name, code, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)`

	level.PrepareCreate(nil)

	_, err := r.db.ExecContext(ctx, query,
		level.ID,
		level.OrganizationID,
		level.Name,
		level.Code,
		level.CreatedAt,
		level.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create education level: %w", err)
	}

	return nil
}

func (r *EducationLevelRepoPostgres) GetByID(ctx context.Context, id uuid.UUID) (*domain.EducationLevel, error) {
	query := `
		SELECT id, organization_id, name, code, created_at, updated_at
		FROM education_levels
		WHERE id = $1 AND deleted_at IS NULL`

	level := &domain.EducationLevel{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&level.ID,
		&level.OrganizationID,
		&level.Name,
		&level.Code,
		&level.CreatedAt,
		&level.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get education level by id: %w", err)
	}

	return level, nil
}
