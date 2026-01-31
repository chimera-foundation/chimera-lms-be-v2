package seed

import (
	"context"
	"errors"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/subject/domain"
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/shared/auth"
	"github.com/google/uuid"
)

type SubjectSeeder struct {
	r domain.SubjectRepository
}

func NewSubjectSeeder(r domain.SubjectRepository) *SubjectSeeder {
	return &SubjectSeeder{
		r: r,
	}
}

func (s *SubjectSeeder) SeedSubjects(ctx context.Context, educationLevelID uuid.UUID) ([]*domain.Subject, error) {
	orgID, ok := auth.GetOrgID(ctx)
	if !ok {
		return nil, errors.New("Organization not found")
	}

	subjects := []*domain.Subject{
		{
			OrganizationID: orgID,
			EducationLevelID: educationLevelID,
			Name: "Calculus",
			Code: "MTK-SMA",
		},
		{
			OrganizationID: orgID,
			EducationLevelID: educationLevelID,
			Name: "Biology",
			Code: "BIO-SMA",
		},
	}

	for _, subject := range subjects {
		err := s.r.Create(ctx, subject)
		if err != nil {
			return nil, err
		}
	}

	return subjects, nil
}
