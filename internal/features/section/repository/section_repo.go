package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/section/domain"
	"github.com/google/uuid"
)

type SectionRepositoryPostgres struct {
	db *sql.DB
}

func NewSectionRepository(db *sql.DB) domain.SectionRepository {
	return & SectionRepositoryPostgres {
		db: db,
	}
}

func (r *SectionRepositoryPostgres) GetSectionIDsByUserID(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	query := `
			SELECT section_id FROM section_members
			WHERE user_id = $1`
	
	
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

func (r *SectionRepositoryPostgres) GetByID(ctx context.Context, id uuid.UUID) (*domain.Section, error) {
	query := `
			SELECT id, cohort_id, name, room_code, capacity FROM sections
			WHERE id = $1`

	section := &domain.Section{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&section.ID,
		&section.CohortID,
		&section.Name,
		&section.RoomCode,
		&section.Capacity,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get section by id: %w", err)
	}

	return section, nil
}