package seed

import (
	"context"
	"errors"
	"fmt"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/course/domain"
	subDomain "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/subject/domain"
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/shared/auth"
	"github.com/google/uuid"
)

type CourseSeeder struct {
	r domain.CourseRepository
}

func NewCourseSeeder(r domain.CourseRepository) *CourseSeeder {
	return &CourseSeeder{
		r: r,
	}
}

func (s *CourseSeeder) SeedCourses(ctx context.Context, subjects []*subDomain.Subject, teacherID uuid.UUID, educationLevelID uuid.UUID) ([]*domain.Course, error) {
	orgID, ok := auth.GetOrgID(ctx)
	if !ok {
		return nil, errors.New("Organization ID doesn't exist")
	}

	var courses []*domain.Course

	for _, subject := range subjects {
		// Grade 10 Course
		course10 := &domain.Course{
			OrganizationID:   orgID,
			InstructorID:     teacherID,
			SubjectID:        subject.ID,
			EducationLevelID: educationLevelID,
			Title:            fmt.Sprintf("%s Kelas 10", subject.Name),
			Description:      fmt.Sprintf("Mata Pelajaran %s untuk Kelas 10", subject.Name),
			Status:           domain.CourseStatus("published"),
			Price:            1000000,
			GradeLevel:       10,
		}
		courses = append(courses, course10)

		// Grade 11 Course
		course11 := &domain.Course{
			OrganizationID:   orgID,
			InstructorID:     teacherID,
			SubjectID:        subject.ID,
			EducationLevelID: educationLevelID,
			Title:            fmt.Sprintf("%s Kelas 11", subject.Name),
			Description:      fmt.Sprintf("Mata Pelajaran %s untuk Kelas 11", subject.Name),
			Status:           domain.CourseStatus("published"),
			Price:            1000000,
			GradeLevel:       11,
		}
		courses = append(courses, course11)
	}

	for _, c := range courses {
		err := s.r.Create(ctx, c)
		if err != nil {
			return nil, err
		}
	}

	return courses, nil
}
