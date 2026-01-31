package seed

import (
	"context"
	"fmt"

	courseDomain "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/course/domain"
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/program/domain"
)

type ProgramCourseSeeder struct {
	r domain.ProgramCourseRepository
}

func NewProgramCourseSeeder(r domain.ProgramCourseRepository) *ProgramCourseSeeder {
	return &ProgramCourseSeeder{
		r: r,
	}
}

func (s *ProgramCourseSeeder) SeedProgramCourses(
	ctx context.Context,
	programs []*domain.Program,
	courses []*courseDomain.Course,
) ([]*domain.ProgramCourse, error) {
	var seededProgramCourses []*domain.ProgramCourse

	// Find MIPA program
	var mipaProgram *domain.Program
	for _, p := range programs {
		if p.Name == "Matematika dan Ilmu Pengetahuan Alam (MIPA)" {
			mipaProgram = p
			break
		}
	}

	if mipaProgram == nil {
		return nil, fmt.Errorf("MIPA program not found")
	}

	// Assign all courses to MIPA program (since we have Calculus and Biology)
	for i, course := range courses {
		pc := &domain.ProgramCourse{
			ProgramID:  mipaProgram.ID,
			CourseID:   course.ID,
			OrderIndex: i + 1,
		}
		if err := s.r.Create(ctx, pc); err != nil {
			return nil, fmt.Errorf("failed to assign course %s to program: %w", course.Title, err)
		}
		seededProgramCourses = append(seededProgramCourses, pc)
	}

	return seededProgramCourses, nil
}
