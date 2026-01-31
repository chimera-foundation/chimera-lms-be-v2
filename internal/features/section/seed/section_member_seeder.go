package seed

import (
	"context"
	"fmt"

	sectionDomain "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/section/domain"
	userDomain "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/user/domain"
)

type SectionMemberSeeder struct {
	r sectionDomain.SectionMemberRepository
}

func NewSectionMemberSeeder(r sectionDomain.SectionMemberRepository) *SectionMemberSeeder {
	return &SectionMemberSeeder{
		r: r,
	}
}

func (s *SectionMemberSeeder) SeedSectionMembers(
	ctx context.Context,
	seededUsers map[string]*userDomain.User,
	sections []*sectionDomain.Section,
) ([]*sectionDomain.SectionMember, error) {
	var seededMembers []*sectionDomain.SectionMember

	// Find sections by name
	var section10, section11 *sectionDomain.Section
	for _, sec := range sections {
		switch sec.Name {
		case "10-A":
			section10 = sec
		case "11-A":
			section11 = sec
		}
	}

	// Assign student10 to Section 10-A as student
	student10 := seededUsers["student10@candletree.com"]
	if student10 != nil && section10 != nil {
		member := &sectionDomain.SectionMember{
			SectionID: section10.ID,
			UserID:    student10.ID,
			RoleType:  sectionDomain.Student,
		}
		if err := s.r.Create(ctx, member); err != nil {
			return nil, fmt.Errorf("failed to add student10 to section: %w", err)
		}
		seededMembers = append(seededMembers, member)
	}

	// Assign student11 to Section 11-A as student
	student11 := seededUsers["student11@candletree.com"]
	if student11 != nil && section11 != nil {
		member := &sectionDomain.SectionMember{
			SectionID: section11.ID,
			UserID:    student11.ID,
			RoleType:  sectionDomain.Student,
		}
		if err := s.r.Create(ctx, member); err != nil {
			return nil, fmt.Errorf("failed to add student11 to section: %w", err)
		}
		seededMembers = append(seededMembers, member)
	}

	// Assign teacher to both sections as teacher
	teacher := seededUsers["teacher@candletree.com"]
	if teacher != nil {
		if section10 != nil {
			member := &sectionDomain.SectionMember{
				SectionID: section10.ID,
				UserID:    teacher.ID,
				RoleType:  sectionDomain.Teacher,
			}
			if err := s.r.Create(ctx, member); err != nil {
				return nil, fmt.Errorf("failed to add teacher to section 10: %w", err)
			}
			seededMembers = append(seededMembers, member)
		}
		if section11 != nil {
			member := &sectionDomain.SectionMember{
				SectionID: section11.ID,
				UserID:    teacher.ID,
				RoleType:  sectionDomain.Teacher,
			}
			if err := s.r.Create(ctx, member); err != nil {
				return nil, fmt.Errorf("failed to add teacher to section 11: %w", err)
			}
			seededMembers = append(seededMembers, member)
		}
	}

	return seededMembers, nil
}
