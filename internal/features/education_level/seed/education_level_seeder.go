package seed

import (
	"context"

	"errors"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/education_level/domain"
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/shared/auth"
)

type EducationLevelSeeder struct {
	elr domain.EducationLevelRepository
}

func NewEducationLevelRepository(elr domain.EducationLevelRepository) *EducationLevelSeeder {
	return &EducationLevelSeeder{
		elr: elr,
	}
}

func (s* EducationLevelSeeder) SeedEducationLevels(ctx context.Context) error {
	orgID, ok := auth.GetOrgID(ctx)
	if !ok {
		return errors.New("Organization ID not found")
	}

	educationLevels := []*domain.EducationLevel{
		{
			OrganizationID: orgID,
			Name: "High School",
			Code: "HIGH",
		},
	}

	for _, eduLevel := range educationLevels {
		err := s.elr.Create(ctx, eduLevel)
		if err != nil {
			return err
		}
	}

	return nil
}