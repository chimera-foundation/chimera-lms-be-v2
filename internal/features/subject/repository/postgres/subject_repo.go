package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/subject/domain"
	"github.com/google/uuid"
)

type SubjectRepoPostgres struct {
	db *sql.DB
}

func NewSubjectRepository(db *sql.DB) domain.SubjectRepository {
	return &SubjectRepoPostgres{db: db}
}

func (r *SubjectRepoPostgres) Create(ctx context.Context, subject *domain.Subject) error {
	query := `
		INSERT INTO subjects (id, organization_id, education_level_id, name, code, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	subject.PrepareCreate(nil)

	_, err := r.db.ExecContext(ctx, query,
		subject.ID,
		subject.OrganizationID,
		subject.EducationLevelID,
		subject.Name,
		subject.Code,
		subject.CreatedAt,
		subject.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create subject: %w", err)
	}

	return nil
}

func (r *SubjectRepoPostgres) GetByID(ctx context.Context, id uuid.UUID) (*domain.Subject, error) {
	query := `
		SELECT id, organization_id, education_level_id, name, code, created_at, updated_at
		FROM subjects
		WHERE id = $1 AND deleted_at IS NULL`

	subject := &domain.Subject{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&subject.ID,
		&subject.OrganizationID,
		&subject.EducationLevelID,
		&subject.Name,
		&subject.Code,
		&subject.CreatedAt,
		&subject.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get subject by id: %w", err)
	}

	return subject, nil
}
