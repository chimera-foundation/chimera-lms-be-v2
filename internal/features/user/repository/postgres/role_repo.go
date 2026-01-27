package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/user/domain"
	"github.com/google/uuid"
)

type RoleRepoPostgres struct {
	db *sql.DB
}

func NewRoleRepository(db *sql.DB) domain.RoleRepository {
	return &RoleRepoPostgres{db: db}
}

func (r *RoleRepoPostgres) Create(ctx context.Context, role *domain.Role) error {
	query := `
		INSERT INTO roles (id, name, created_at, updated_at)
		VALUES ($1, $2, $3, $4)`

	role.PrepareCreate(nil)

	_, err := r.db.ExecContext(ctx, query,
		role.ID,
		role.Name,
		role.CreatedAt,
		role.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create role: %w", err)
	}

	return nil
}

func (r *RoleRepoPostgres) GetByName(ctx context.Context, name string) (*domain.Role, error) {
	query := `
		SELECT id, name, created_at, updated_at
		FROM roles
		WHERE name = $1 AND deleted_at IS NULL`

	role := &domain.Role{}
	err := r.db.QueryRowContext(ctx, query, name).Scan(
		&role.ID,
		&role.Name,
		&role.CreatedAt,
		&role.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get role by name: %w", err)
	}

	return role, nil
}

func (r *RoleRepoPostgres) AssignRoleToUser(ctx context.Context, userID, roleID uuid.UUID) error {
	query := `INSERT INTO user_roles (user_id, role_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`

	_, err := r.db.ExecContext(ctx, query, userID, roleID)
	if err != nil {
		return fmt.Errorf("failed to assign role to user: %w", err)
	}

	return nil
}
