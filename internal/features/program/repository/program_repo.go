package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/program/domain"
	"github.com/google/uuid"
)

type ProgramRepoPostgres struct {
	db *sql.DB
}

func NewProgramRepository(db *sql.DB) domain.ProgramRepository {
	return &ProgramRepoPostgres{db: db}
}

func (r *ProgramRepoPostgres) Create(ctx context.Context, program *domain.Program) error {
	query := `
		INSERT INTO programs (id, organization_id, name, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)`

	program.PrepareCreate(nil)

	_, err := r.db.ExecContext(ctx, query,
		program.ID,
		program.OrganizationID,
		program.Name,
		program.Description,
		program.CreatedAt,
		program.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create program: %w", err)
	}

	return nil
}

func (r *ProgramRepoPostgres) GetByID(ctx context.Context, id uuid.UUID) (*domain.Program, error) {
	query := `
		SELECT id, organization_id, name, description, created_at, updated_at
		FROM programs
		WHERE id = $1 AND deleted_at IS NULL`

	program := &domain.Program{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&program.ID,
		&program.OrganizationID,
		&program.Name,
		&program.Description,
		&program.CreatedAt,
		&program.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get program by id: %w", err)
	}

	return program, nil
}
