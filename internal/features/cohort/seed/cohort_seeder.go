package seed

import (
	"context"
	"errors"
	"fmt"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/cohort/domain"
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/shared/auth"
	"github.com/google/uuid"
)

type CohortSeeder struct {
	r domain.CohortRepository
}

func NewCohortSeeder(r domain.CohortRepository) *CohortSeeder {
	return &CohortSeeder{
		r: r,
	}
}

func (s *CohortSeeder) SeedCohorts(ctx context.Context, academicPeriodID uuid.UUID, educationLevelID uuid.UUID) ([]*domain.Cohort, error) {
	orgID, ok := auth.GetOrgID(ctx)
	if !ok {
		return nil, errors.New("Organization ID doesn't exist")
	}

	cohorts := []*domain.Cohort{
		{
			OrganizationID:   orgID,
			AcademicPeriodID: academicPeriodID,
			EducationLevelID: educationLevelID,
			Name:             "Grade 10 Cohort",
		},
		{
			OrganizationID:   orgID,
			AcademicPeriodID: academicPeriodID,
			EducationLevelID: educationLevelID,
			Name:             "Grade 11 Cohort",
		},
	}

	var seededCohorts []*domain.Cohort

	for _, c := range cohorts {
		if err := s.r.Create(ctx, c); err != nil {
			return nil, fmt.Errorf("failed to create cohort %s: %w", c.Name, err)
		}
		seededCohorts = append(seededCohorts, c)
	}

	return seededCohorts, nil
}
