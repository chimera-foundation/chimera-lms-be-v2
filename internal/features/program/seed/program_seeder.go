package seed

import (
	"context"
	"errors"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/program/domain"
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/shared/auth"
)

type ProgramSeeder struct {
	r domain.ProgramRepository
}

func NewProgramSeeder(r domain.ProgramRepository) *ProgramSeeder {
	return &ProgramSeeder{
		r: r,
	}
}

func (s *ProgramSeeder) SeedPrograms(ctx context.Context) error {
	orgID, ok := auth.GetOrgID(ctx)
	if !ok {
		return errors.New("Organization ID not found") 
	}

	programs := []domain.Program {
		{
			OrganizationID: orgID,
			Name: "Matematika dan Ilmu Pengetahuan Alam (MIPA)",
			Description: "Sesuai judul bg",
		},
		{
			OrganizationID: orgID,
			Name: "Ilmu Pengetahuan Sosial (IPS)",
			Description: "Sesuai sejarah bg",
		},
	}

	for i := range programs {
		err := s.r.Create(ctx, &programs[i])
		if err != nil {
			return errors.New("Failed at creating a program")
		}
	}

	return nil
}