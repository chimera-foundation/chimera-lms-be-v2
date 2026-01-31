package seed

import (
	"context"
	"fmt"

	cohortDomain "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/cohort/domain"
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/section/domain"
)

type SectionSeeder struct {
	r domain.SectionRepository
}

func NewSectionSeeder(r domain.SectionRepository) *SectionSeeder {
	return &SectionSeeder{
		r: r,
	}
}

func (s *SectionSeeder) SeedSections(ctx context.Context, cohorts []*cohortDomain.Cohort) ([]*domain.Section, error) {
	var seededSections []*domain.Section

	for _, cohort := range cohorts {
		var sectionName string
		var roomCode string
		switch cohort.Name {
		case "Grade 10 Cohort":
			sectionName = "10-A"
			roomCode=  "R-10A"
		case "Grade 11 Cohort":
			sectionName = "11-A"
			roomCode=  "R-11A"
		default:
			sectionName = "Unknown Section"
		}

		section := &domain.Section{
			CohortID: cohort.ID,
			Name:     sectionName,
			RoomCode: roomCode,
			Capacity: 30,
		}

		if err := s.r.Create(ctx, section); err != nil {
			return nil, fmt.Errorf("failed to create section %s for cohort %s: %w", sectionName, cohort.Name, err)
		}
		seededSections = append(seededSections, section)
	}

	return seededSections, nil
}
