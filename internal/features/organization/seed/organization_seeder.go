package seed

import (
	"context"
	"time"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/organization/domain"
)

type OrganizationSeeder struct {
	or  domain.OrganizationRepository
	apr domain.AcademicPeriodRepository
}

func NewOrganizationSeeder(or domain.OrganizationRepository, apr domain.AcademicPeriodRepository) *OrganizationSeeder {
	return &OrganizationSeeder{
		or:  or,
		apr: apr,
	}
}

func (s *OrganizationSeeder) SeedOrganizations(ctx context.Context) (*domain.Organization, *domain.AcademicPeriod, error) {
	mock_period := domain.NewAcademicPeriod(
		"2026/2027 genap",
		time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC), // Typical start of Ganjil
		time.Date(2026, time.June, 31, 23, 59, 59, 0, time.UTC),
	)

	org_type := domain.OrgType("high_school")

	organization := domain.NewOrganization(
		"Candle Tree School",
		"cts",
		org_type,
		"Mars, Planet ke Gaktau",
		nil,
	)

	organization.PrepareCreate(nil)

	err := s.or.Create(ctx, organization)
	if err != nil {
		return nil, nil, err
	}

	err = s.apr.Create(ctx, mock_period, organization.ID)
	if err != nil {
		return nil, nil, err
	}

	return organization, mock_period, nil
}
