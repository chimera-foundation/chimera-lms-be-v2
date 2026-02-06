package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/organization/domain"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type OrganizationRepoPostgres struct {
	db  *sql.DB
	log *logrus.Logger
}

func NewOrganizationRepo(db *sql.DB, log *logrus.Logger) domain.OrganizationRepository {
	return &OrganizationRepoPostgres{db: db, log: log}
}

func (r *OrganizationRepoPostgres) Create(ctx context.Context, org *domain.Organization) error {
	query := `
		INSERT INTO organizations (id, name, slug, type, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)`
	org.PrepareCreate(nil)
	_, err := r.db.ExecContext(ctx, query,
		org.ID,
		org.Name,
		org.Slug,
		org.Type,
		org.CreatedAt,
		org.UpdatedAt,
	)
	if err != nil {
		r.log.WithError(err).WithField("org_id", org.ID).Error("failed to insert organization")
		return fmt.Errorf("failed to insert organization: %w", err)
	}
	r.log.WithFields(logrus.Fields{"org_id": org.ID, "name": org.Name}).Info("organization created successfully")
	return nil
}

func (r *OrganizationRepoPostgres) GetByID(ctx context.Context, id uuid.UUID) (*domain.Organization, error) {
	query := `
			SELECT id, name, slug, type, created_at, updated_at, is_system_org 
			FROM organizations 
			WHERE id = $1 AND deleted_at IS NULL`

	organization := &domain.Organization{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&organization.ID,
		&organization.Name,
		&organization.Slug,
		&organization.Type,
		&organization.CreatedAt,
		&organization.UpdatedAt,
		&organization.IsSystemOrg,
	)

	if err == sql.ErrNoRows {
		r.log.WithField("org_id", id).Debug("organization not found by id")
		return nil, nil
	}
	if err != nil {
		r.log.WithError(err).WithField("org_id", id).Error("failed to get organization by id")
		return nil, fmt.Errorf("failed to get organization by id: %w", err)
	}

	return organization, nil
}

func (r *OrganizationRepoPostgres) Update(ctx context.Context, org *domain.Organization) error {
	query := `
		UPDATE organizations
		SET name=$2, slug=$3, type=$4, created_at=$5, updated_at=$6
		WHERE id=$1 AND deleted_at IS NULL`

	res, err := r.db.ExecContext(ctx, query,
		org.ID,
		org.Name,
		org.Slug,
		org.Type,
		org.CreatedAt,
		org.UpdatedAt,
	)

	if err != nil {
		r.log.WithError(err).WithField("org_id", org.ID).Error("failed to update organization")
		return fmt.Errorf("failed to update organization: %w", err)
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		r.log.WithField("org_id", org.ID).Warn("organization not found or already deleted for update")
		return fmt.Errorf("organization not found or already deleted")
	}

	r.log.WithField("org_id", org.ID).Info("organization updated successfully")
	return nil
}

func (r *OrganizationRepoPostgres) Delete(ctx context.Context, orgID uuid.UUID) error {
	query := `
		DELETE FROM organizations
		WHERE id=$1`

	res, err := r.db.ExecContext(ctx, query, orgID)

	if err != nil {
		r.log.WithError(err).WithField("org_id", orgID).Error("failed to delete organization")
		return fmt.Errorf("failed to delete organization: %w", err)
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		r.log.WithField("org_id", orgID).Warn("organization not found or already deleted for deletion")
		return fmt.Errorf("organization not found or already deleted")
	}

	r.log.WithField("org_id", orgID).Info("organization deleted successfully")
	return nil
}

func (r *OrganizationRepoPostgres) GetBySlug(ctx context.Context, slug string) (*domain.Organization, error) {
	query := `
		SELECT id, name, slug, type, created_at, updated_at, is_system_org
		FROM organizations
		WHERE slug = $1 AND deleted_at IS NULL`

	organization := &domain.Organization{}
	err := r.db.QueryRowContext(ctx, query, slug).Scan(
		&organization.ID,
		&organization.Name,
		&organization.Slug,
		&organization.Type,
		&organization.CreatedAt,
		&organization.UpdatedAt,
		&organization.IsSystemOrg,
	)

	if err == sql.ErrNoRows {
		r.log.WithField("slug", slug).Debug("organization not found by slug")
		return nil, nil
	}
	if err != nil {
		r.log.WithError(err).WithField("slug", slug).Error("failed to get organization by slug")
		return nil, fmt.Errorf("failed to get organization by slug: %w", err)
	}

	return organization, nil
}

func (r *OrganizationRepoPostgres) GetIDByUserID(ctx context.Context, userID uuid.UUID) (uuid.UUID, error) {
	query := `
        SELECT organization_id 
        FROM users 
        WHERE id = $1 AND deleted_at IS NULL`

	var orgID uuid.UUID
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&orgID)

	if err != nil {
		if err == sql.ErrNoRows {
			r.log.WithField("user_id", userID).Debug("organization not found for this user")
			return uuid.Nil, errors.New("organization not found for this user")
		}
		r.log.WithError(err).WithField("user_id", userID).Error("failed to get org id by user id")
		return uuid.Nil, err
	}

	return orgID, nil
}
