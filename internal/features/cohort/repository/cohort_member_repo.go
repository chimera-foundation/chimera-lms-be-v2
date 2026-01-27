package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/cohort/domain"
)

type CohortMemberRepoPostgres struct {
	db *sql.DB
}

func NewCohortMemberRepository(db *sql.DB) domain.CohortMemberRepository {
	return &CohortMemberRepoPostgres{db: db}
}

func (r *CohortMemberRepoPostgres) Create(ctx context.Context, member *domain.CohortMember) error {
	query := `
		INSERT INTO cohort_members (cohort_id, user_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT DO NOTHING`

	now := time.Now()
	member.CreatedAt = now
	member.UpdatedAt = now

	_, err := r.db.ExecContext(ctx, query,
		member.CohortID,
		member.UserID,
		member.CreatedAt,
		member.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create cohort member: %w", err)
	}

	return nil
}
