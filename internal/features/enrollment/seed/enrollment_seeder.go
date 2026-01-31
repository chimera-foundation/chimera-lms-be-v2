package seed

import (
	"context"
	"fmt"
	"time"

	courseDomain "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/course/domain"
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/enrollment/domain"
	sectionDomain "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/section/domain"
	userDomain "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/user/domain"
	"github.com/google/uuid"
)

type EnrollmentSeeder struct {
	r domain.EnrollmentRepository
}

func NewEnrollmentSeeder(r domain.EnrollmentRepository) *EnrollmentSeeder {
	return &EnrollmentSeeder{
		r: r,
	}
}

func (s *EnrollmentSeeder) SeedEnrollments(
	ctx context.Context,
	seededUsers map[string]*userDomain.User,
	courses []*courseDomain.Course,
	sections []*sectionDomain.Section,
	academicPeriodID uuid.UUID,
) ([]*domain.Enrollment, error) {

	var seededEnrollments []*domain.Enrollment

	// Helper to find section by name
	findSection := func(name string) *sectionDomain.Section {
		for _, s := range sections {
			if s.Name == name {
				return s
			}
		}
		return nil
	}

	// 1. Enroll Student 10 to Grade 10 Courses
	student10 := seededUsers["student10@candletree.com"]
	section10 := findSection("10-A")

	for _, course := range courses {
		if course.GradeLevel == 10 && student10 != nil && section10 != nil {
			enrollment := &domain.Enrollment{
				UserID:           student10.ID,
				CourseID:         course.ID,
				SectionID:        section10.ID,
				AcademicPeriodID: academicPeriodID,
				Status:           domain.Active,
				EnrolledAt:       time.Now(),
			}

			if err := s.r.Create(ctx, enrollment); err != nil {
				return nil, fmt.Errorf("failed to enroll student 10 to course %s: %w", course.Title, err)
			}
			seededEnrollments = append(seededEnrollments, enrollment)
		}
	}

	// 2. Enroll Student 11 to Grade 11 Courses
	student11 := seededUsers["student11@candletree.com"]
	section11 := findSection("11-A")

	for _, course := range courses {
		if course.GradeLevel == 11 && student11 != nil && section11 != nil {
			enrollment := &domain.Enrollment{
				UserID:           student11.ID,
				CourseID:         course.ID,
				SectionID:        section11.ID,
				AcademicPeriodID: academicPeriodID,
				Status:           domain.Active,
				EnrolledAt:       time.Now(),
			}

			if err := s.r.Create(ctx, enrollment); err != nil {
				return nil, fmt.Errorf("failed to enroll student 11 to course %s: %w", course.Title, err)
			}
			seededEnrollments = append(seededEnrollments, enrollment)
		}
	}

	return seededEnrollments, nil
}
