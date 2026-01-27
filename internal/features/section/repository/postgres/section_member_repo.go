package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/section/domain"
)

type SectionMemberRepoPostgres struct {
	db *sql.DB
}

func NewSectionMemberRepository(db *sql.DB) domain.SectionMemberRepository {
	return &SectionMemberRepoPostgres{db: db}
}

func (r *SectionMemberRepoPostgres) Create(ctx context.Context, member *domain.SectionMember) error {
	query := `
		INSERT INTO section_members (section_id, user_id, role_type)
		VALUES ($1, $2, $3)
		ON CONFLICT DO NOTHING`

	_, err := r.db.ExecContext(ctx, query,
		member.SectionID,
		member.UserID,
		member.RoleType,
	)

	if err != nil {
		return fmt.Errorf("failed to create section member: %w", err)
	}

	return nil
}
