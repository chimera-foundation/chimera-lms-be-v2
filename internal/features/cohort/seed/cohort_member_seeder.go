package seed

import (
	"context"
	"fmt"

	cohortDomain "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/cohort/domain"
	userDomain "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/user/domain"
)

type CohortMemberSeeder struct {
	r cohortDomain.CohortMemberRepository
}

func NewCohortMemberSeeder(r cohortDomain.CohortMemberRepository) *CohortMemberSeeder {
	return &CohortMemberSeeder{
		r: r,
	}
}

func (s *CohortMemberSeeder) SeedCohortMembers(
	ctx context.Context,
	seededUsers map[string]*userDomain.User,
	cohorts []*cohortDomain.Cohort,
) ([]*cohortDomain.CohortMember, error) {
	var seededMembers []*cohortDomain.CohortMember

	// Find cohorts by name
	var cohort10, cohort11 *cohortDomain.Cohort
	for _, c := range cohorts {
		switch c.Name {
		case "Grade 10 Cohort":
			cohort10 = c
		case "Grade 11 Cohort":
			cohort11 = c
		}
	}

	// Assign student10 to Grade 10 Cohort
	student10 := seededUsers["student10@candletree.com"]
	if student10 != nil && cohort10 != nil {
		member := &cohortDomain.CohortMember{
			CohortID: cohort10.ID,
			UserID:   student10.ID,
		}
		if err := s.r.Create(ctx, member); err != nil {
			return nil, fmt.Errorf("failed to add student10 to cohort: %w", err)
		}
		seededMembers = append(seededMembers, member)
	}

	// Assign student11 to Grade 11 Cohort
	student11 := seededUsers["student11@candletree.com"]
	if student11 != nil && cohort11 != nil {
		member := &cohortDomain.CohortMember{
			CohortID: cohort11.ID,
			UserID:   student11.ID,
		}
		if err := s.r.Create(ctx, member); err != nil {
			return nil, fmt.Errorf("failed to add student11 to cohort: %w", err)
		}
		seededMembers = append(seededMembers, member)
	}

	// Assign teacher to both cohorts
	teacher := seededUsers["teacher@candletree.com"]
	if teacher != nil {
		if cohort10 != nil {
			member := &cohortDomain.CohortMember{
				CohortID: cohort10.ID,
				UserID:   teacher.ID,
			}
			if err := s.r.Create(ctx, member); err != nil {
				return nil, fmt.Errorf("failed to add teacher to cohort 10: %w", err)
			}
			seededMembers = append(seededMembers, member)
		}
		if cohort11 != nil {
			member := &cohortDomain.CohortMember{
				CohortID: cohort11.ID,
				UserID:   teacher.ID,
			}
			if err := s.r.Create(ctx, member); err != nil {
				return nil, fmt.Errorf("failed to add teacher to cohort 11: %w", err)
			}
			seededMembers = append(seededMembers, member)
		}
	}

	return seededMembers, nil
}
