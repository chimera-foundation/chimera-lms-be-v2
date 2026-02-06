package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/section/domain"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type SectionRepositoryPostgres struct {
	db  *sql.DB
	log *logrus.Logger
}

func NewSectionRepository(db *sql.DB, log *logrus.Logger) domain.SectionRepository {
	return &SectionRepositoryPostgres{
		db:  db,
		log: log,
	}
}

func (r *SectionRepositoryPostgres) Create(ctx context.Context, section *domain.Section) error {
	query := `
		INSERT INTO sections (id, cohort_id, name, room_number, capacity, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	section.PrepareCreate(nil)

	_, err := r.db.ExecContext(ctx, query,
		section.ID,
		section.CohortID,
		section.Name,
		section.RoomCode,
		section.Capacity,
		section.CreatedAt,
		section.UpdatedAt,
	)

	if err != nil {
		r.log.WithError(err).WithField("section_id", section.ID).Error("failed to create section")
		return fmt.Errorf("failed to create section: %w", err)
	}

	r.log.WithFields(logrus.Fields{"section_id": section.ID, "name": section.Name}).Info("section created successfully")
	return nil
}

func (r *SectionRepositoryPostgres) GetSectionIDsByUserID(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	query := `
			SELECT section_id FROM section_members
			WHERE user_id = $1`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		r.log.WithError(err).WithField("user_id", userID).Error("failed to query section ids by user")
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
		r.log.WithError(err).WithField("user_id", userID).Error("error iterating section ids")
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
		r.log.WithField("section_id", id).Debug("section not found by id")
		return nil, nil
	}

	if err != nil {
		r.log.WithError(err).WithField("section_id", id).Error("failed to get section by id")
		return nil, fmt.Errorf("failed to get section by id: %w", err)
	}

	return section, nil
}
