package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/course/domain"
	"github.com/google/uuid"
)

type ModuleRepoPostgres struct {
	db *sql.DB
}

func NewModuleRepository(db *sql.DB) domain.ModuleRepository {
	return &ModuleRepoPostgres{db: db}
}

func (r *ModuleRepoPostgres) Create(ctx context.Context, module *domain.Module) error {
	query := `
		INSERT INTO modules (id, course_id, title, order_index, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)`

	module.PrepareCreate(nil)

	_, err := r.db.ExecContext(ctx, query,
		module.ID,
		module.CourseID,
		module.Title,
		module.OrderIndex,
		module.CreatedAt,
		module.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create module: %w", err)
	}

	return nil
}

func (r *ModuleRepoPostgres) GetByID(ctx context.Context, id uuid.UUID) (*domain.Module, error) {
	query := `
		SELECT id, course_id, title, order_index, created_at, updated_at
		FROM modules
		WHERE id = $1 AND deleted_at IS NULL`

	module := &domain.Module{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&module.ID,
		&module.CourseID,
		&module.Title,
		&module.OrderIndex,
		&module.CreatedAt,
		&module.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get module by id: %w", err)
	}

	return module, nil
}
