package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/user/domain"
	"github.com/google/uuid"
)

type UserRepoPostgres struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) domain.UserRepository {
	return &UserRepoPostgres{db: db}
}

func (r *UserRepoPostgres) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (id, organization_id, email, password_hash, is_superuser, created_at, updated_at, first_name, last_name)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	user.PrepareCreate(&user.OrganizationID)
	_, err := r.db.ExecContext(ctx, query,
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
	return nil
}

func (r *UserRepoPostgres) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, organization_id, email, password_hash, is_superuser, created_at, updated_at, first_name, last_name
		FROM users 
		WHERE email = $1 AND deleted_at IS NULL`

	user := &domain.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.OrganizationID,
		&user.Email,
		&user.PasswordHash,
		&user.IsSuperuser,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.FirstName,
		&user.LastName,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return user, nil
}

func (r *UserRepoPostgres) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	query := `
		SELECT id, organization_id, email, password_hash, is_superuser, created_at, updated_at, first_name, last_name
		FROM users 
		WHERE id = $1 AND deleted_at IS NULL`

	user := &domain.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.OrganizationID,
		&user.Email,
		&user.PasswordHash,
		&user.IsSuperuser,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.FirstName,
		&user.LastName,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	return user, nil
}

func (r *UserRepoPostgres) Update(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users 
		SET organization_id = $2, email = $3, password_hash = $4, is_superuser = $5, updated_at = $6
		WHERE id = $1 AND deleted_at IS NULL`

	user.UpdatedAt = time.Now()
	user.UpdatedBy = &user.ID
	res, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.OrganizationID,
		user.Email,
		user.PasswordHash,
		user.IsSuperuser,
		time.Now(),
	)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("user not found or already deleted")
	}

	return nil
}
