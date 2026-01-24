package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/organization/domain"
	"github.com/google/uuid"
)

type OrganizationRepoPostgres struct {
	db *sql.DB
}

func NewOrganizationRepo(db *sql.DB) domain.OrganizationRepository{
	return &OrganizationRepoPostgres{db: db}
}

func (r OrganizationRepoPostgres) Create(ctx context.Context, org *domain.Organization) error {
	query := `
		INSERT INTO organizations (id, name, slug, type, created_at, updated_at)
		VALUE ($1, $2, $3, $4, $5, $6)`
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
		return fmt.Errorf("failed to insert organization: %w", err)
	}
	return nil
}

func (r OrganizationRepoPostgres) Update(ctx context.Context, org *domain.Organization) error {
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
		return fmt.Errorf("failed to update organization: %w", err)
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("organization not found or already deleted")
	}

	return nil
}

func (r OrganizationRepoPostgres) Delete(ctx context.Context, orgID uuid.UUID) error {
	query := `
		DELETE FROM organizations
		WHERE id=$1`

	res, err := r.db.Exec(query, orgID)

	if err != nil {
		return fmt.Errorf("failed to delete organization: %w", err)
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("organization not found or already deleted")
	}

	return nil
}

func (r *OrganizationRepoPostgres) GetIDByUserID(ctx context.Context, userID uuid.UUID) (uuid.UUID, error) {
    // We query the 'users' table because it holds the foreign key 
    // to the organization according to your schema.
    query := `
        SELECT organization_id 
        FROM users 
        WHERE id = $1 AND deleted_at IS NULL`
    
    var orgID uuid.UUID
    err := r.db.QueryRowContext(ctx, query, userID).Scan(&orgID)

    if err != nil {
        if err == sql.ErrNoRows {
            return uuid.Nil, errors.New("organization not found for this user")
        }
        return uuid.Nil, err
    }

    return orgID, nil
}