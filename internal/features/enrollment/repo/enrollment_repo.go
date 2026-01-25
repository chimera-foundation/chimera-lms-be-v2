package repo

import (
	"context"
	"database/sql"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/enrollment/domain"
	"github.com/google/uuid"
)

type EnrollmentRepositoryPostgres struct {
	db *sql.DB
}

func NewEnrollmentRepository(db *sql.DB) domain.EnrollmentRepository {
	return &EnrollmentRepositoryPostgres{
		db: db,
	}
}

func (r *EnrollmentRepositoryPostgres) GetActiveSectionIDsByUserID(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	query := `
			SELECT section_id 
			FROM enrollments 
			WHERE user_id = $1 AND status = 'active'`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []uuid.UUID
	for rows.Next() {
		var id uuid.UUID

		if err := rows.Scan(&id); err != nil {
			return nil, err
		}

		result = append(result, id)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
