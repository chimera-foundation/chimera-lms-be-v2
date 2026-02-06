package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/user/domain"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type RoleRepoPostgres struct {
	db  *sql.DB
	log *logrus.Logger
}

func NewRoleRepository(db *sql.DB, log *logrus.Logger) domain.RoleRepository {
	return &RoleRepoPostgres{db: db, log: log}
}

func (r *RoleRepoPostgres) Create(ctx context.Context, role *domain.Role) error {
	query := `
		INSERT INTO roles (id, name, permissions, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)`

	permsJSON, err := json.Marshal(role.Permissions)
	if err != nil {
		r.log.WithError(err).WithField("role_name", role.Name).Error("failed to marshal permissions")
		return fmt.Errorf("failed to marshal permissions: %w", err)
	}
	role.PrepareCreate(nil)

	_, err = r.db.ExecContext(ctx, query,
		role.ID,
		role.Name,
		permsJSON,
		role.CreatedAt,
		role.UpdatedAt,
	)

	if err != nil {
		r.log.WithError(err).WithField("role_name", role.Name).Error("failed to create role")
		return fmt.Errorf("failed to create role: %w", err)
	}

	r.log.WithField("role_id", role.ID).Info("role created successfully")
	return nil
}

func (r *RoleRepoPostgres) GetByName(ctx context.Context, name string) (*domain.Role, error) {
	query := `
        SELECT id, name, permissions, created_at, updated_at
        FROM roles
        WHERE name = $1 AND deleted_at IS NULL`

	role := &domain.Role{}
	var permissionsJSON []byte

	err := r.db.QueryRowContext(ctx, query, name).Scan(
		&role.ID,
		&role.Name,
		&permissionsJSON,
		&role.CreatedAt,
		&role.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		r.log.WithField("role_name", name).Debug("role not found by name")
		return nil, nil
	}
	if err != nil {
		r.log.WithError(err).WithField("role_name", name).Error("failed to get role")
		return nil, fmt.Errorf("failed to get role: %w", err)
	}

	if len(permissionsJSON) > 0 {
		if err := json.Unmarshal(permissionsJSON, &role.Permissions); err != nil {
			return nil, fmt.Errorf("failed to unmarshal permissions: %w", err)
		}
	}

	return role, nil
}

func (r *RoleRepoPostgres) AssignRoleToUser(ctx context.Context, userID, roleID uuid.UUID) error {
	query := `INSERT INTO user_roles (user_id, role_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`

	_, err := r.db.ExecContext(ctx, query, userID, roleID)
	if err != nil {
		r.log.WithError(err).WithFields(logrus.Fields{"user_id": userID, "role_id": roleID}).Error("failed to assign role to user")
		return fmt.Errorf("failed to assign role to user: %w", err)
	}

	r.log.WithFields(logrus.Fields{"user_id": userID, "role_id": roleID}).Info("role assigned to user")
	return nil
}

func (r *RoleRepoPostgres) RevokeUserRole(ctx context.Context, userID, roleID uuid.UUID) error {
	query := `DELETE FROM user_roles WHERE user_id = $1 AND role_id = $2`

	_, err := r.db.ExecContext(ctx, query, userID, roleID)
	if err != nil {
		r.log.WithError(err).WithFields(logrus.Fields{"user_id": userID, "role_id": roleID}).Error("failed to revoke role")
		return fmt.Errorf("failed to revoke role: %w", err)
	}

	r.log.WithFields(logrus.Fields{"user_id": userID, "role_id": roleID}).Info("role revoked from user")
	return nil
}
