package seed

import (
	"context"
	"fmt"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/course/domain"
)

type ModuleSeeder struct {
	r domain.ModuleRepository
}

func NewModuleSeeder(r domain.ModuleRepository) *ModuleSeeder {
	return &ModuleSeeder{
		r: r,
	}
}

func (s *ModuleSeeder) SeedModules(ctx context.Context, courses []*domain.Course) ([]*domain.Module, error) {
	var seededModules []*domain.Module

	for _, course := range courses {
		// Create 1 Module per Course
		module := &domain.Module{
			CourseID:   course.ID,
			Title:      fmt.Sprintf("Module 1: Introduction to %s", course.Title),
			OrderIndex: 1,
		}

		if err := s.r.Create(ctx, module); err != nil {
			return nil, fmt.Errorf("failed to create module for course %s: %w", course.Title, err)
		}
		seededModules = append(seededModules, module)
	}

	return seededModules, nil
}
