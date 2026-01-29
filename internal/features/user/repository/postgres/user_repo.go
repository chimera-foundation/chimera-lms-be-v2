package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/user/domain"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type UserRepoPostgres struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) domain.UserRepository {
	return &UserRepoPostgres{db: db}
}

func (r *UserRepoPostgres) Create(ctx context.Context, user *domain.User) error {
    tx, err := r.db.BeginTx(ctx, nil)
    if err != nil {
        return fmt.Errorf("failed to begin transaction: %w", err)
    }

    defer tx.Rollback()

    user.PrepareCreate(&user.OrganizationID)

    userQuery := `
        INSERT INTO users (
            id, 
            organization_id, 
            email, 
            password_hash, 
            is_superuser, 
            created_at, 
            updated_at, 
            first_name, 
            last_name
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

    _, err = tx.ExecContext(ctx, userQuery,
        user.ID,
        user.OrganizationID,
        user.Email,
        user.PasswordHash,
        user.IsSuperuser,
        user.CreatedAt,
        user.UpdatedAt,
        user.FirstName,
        user.LastName,
    )
    if err != nil {
        return fmt.Errorf("failed to insert user: %w", err)
    }

    if len(user.Roles) > 0 {
        roleQuery := `
            INSERT INTO user_roles (user_id, role_id)
            SELECT $1, unnest($2::uuid[])`
        
        roleIDs := make([]uuid.UUID, len(user.Roles))
        for i, role := range user.Roles {
            roleIDs[i] = role.ID
        }

        _, err = tx.ExecContext(ctx, roleQuery, user.ID, pq.Array(roleIDs))
        if err != nil {
            return fmt.Errorf("failed to bulk assign roles: %w", err)
        }
    }

    if err := tx.Commit(); err != nil {
        return fmt.Errorf("failed to commit transaction: %w", err)
    }

    return nil
}

func (r *UserRepoPostgres) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
    query := `
        SELECT 
            u.id, u.organization_id, u.email, u.password_hash, u.first_name, u.last_name, 
            u.is_superuser, u.created_at, u.updated_at,
            COALESCE(
                (SELECT jsonb_agg(jsonb_build_object('id', r.id, 'name', r.name))
                 FROM user_roles ur
                 JOIN roles r ON ur.role_id = r.id
                 WHERE ur.user_id = u.id), 
            '[]') as roles_json
        FROM users u
        WHERE u.email = $1 AND u.deleted_at IS NULL`

    user := &domain.User{}
    var rolesJSON []byte

    err := r.db.QueryRowContext(ctx, query, email).Scan(
        &user.ID,
        &user.OrganizationID,
        &user.Email,
        &user.PasswordHash,
        &user.FirstName,
        &user.LastName,
        &user.IsSuperuser,
        &user.CreatedAt,
        &user.UpdatedAt,
        &rolesJSON,
    )

    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, fmt.Errorf("failed to get user by email: %w", err)
    }

    if err := json.Unmarshal(rolesJSON, &user.Roles); err != nil {
        return nil, fmt.Errorf("failed to unmarshal roles: %w", err)
    }

    return user, nil
}

func (r *UserRepoPostgres) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
    query := `
        SELECT 
            u.id, u.organization_id, u.email, u.password_hash, u.first_name, u.last_name, 
            u.is_superuser, u.created_at, u.updated_at,
            COALESCE(
                (SELECT jsonb_agg(jsonb_build_object('id', r.id, 'name', r.name))
                 FROM user_roles ur
                 JOIN roles r ON ur.role_id = r.id
                 WHERE ur.user_id = u.id), 
            '[]') as roles_json
        FROM users u
        WHERE u.id = $1 AND u.deleted_at IS NULL`

    user := &domain.User{}
    var rolesJSON []byte

    err := r.db.QueryRowContext(ctx, query, id).Scan(
        &user.ID,
        &user.OrganizationID,
        &user.Email,
        &user.PasswordHash,
        &user.FirstName,
        &user.LastName,
        &user.IsSuperuser,
        &user.CreatedAt,
        &user.UpdatedAt,
        &rolesJSON,
    )

    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, fmt.Errorf("failed to get user by id: %w", err)
    }

    // Unmarshal the JSON array directly into the User.Roles slice
    if err := json.Unmarshal(rolesJSON, &user.Roles); err != nil {
        return nil, fmt.Errorf("failed to unmarshal roles: %w", err)
    }

    return user, nil
}

func (r *UserRepoPostgres) Update(ctx context.Context, user *domain.User) error {
    tx, err := r.db.BeginTx(ctx, nil)
    if err != nil {
        return fmt.Errorf("begin tx: %w", err)
    }
    defer tx.Rollback()

    user.UpdatedAt = time.Now()
    
    userQuery := `
        UPDATE users 
        SET email = $2, password_hash = $3, first_name = $4, last_name = $5, is_superuser = $6, updated_at = $7
        WHERE id = $1 AND deleted_at IS NULL`

    _, err = tx.ExecContext(ctx, userQuery,
        user.ID, user.Email, user.PasswordHash, user.FirstName, user.LastName, user.IsSuperuser, user.UpdatedAt,
    )
    if err != nil {
        return fmt.Errorf("update user: %w", err)
    }

    _, err = tx.ExecContext(ctx, "DELETE FROM user_roles WHERE user_id = $1", user.ID)
    if err != nil {
        return fmt.Errorf("clear roles: %w", err)
    }

    if len(user.Roles) > 0 {
        roleIDs := make([]uuid.UUID, len(user.Roles))
        for i, role := range user.Roles {
            roleIDs[i] = role.ID
        }

        roleQuery := `INSERT INTO user_roles (user_id, role_id) SELECT $1, unnest($2::uuid[])`
        _, err = tx.ExecContext(ctx, roleQuery, user.ID, pq.Array(roleIDs))
        if err != nil {
            return fmt.Errorf("sync roles: %w", err)
        }
    }

    return tx.Commit()
}
