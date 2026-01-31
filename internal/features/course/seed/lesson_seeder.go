package seed

import (
	"context"
	"fmt"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/course/domain"
)

type LessonSeeder struct {
	r domain.LessonRepository
}

func NewLessonSeeder(r domain.LessonRepository) *LessonSeeder {
	return &LessonSeeder{
		r: r,
	}
}

func (s *LessonSeeder) SeedLessons(ctx context.Context, modules []*domain.Module) ([]*domain.Lesson, error) {
	var seededLessons []*domain.Lesson

	for _, module := range modules {
		for i := 1; i <= 12; i++ {
			lesson := &domain.Lesson{
				ModuleID:   module.ID,
				Title:      fmt.Sprintf("Lesson %d: Topic %d", i, i),
				OrderIndex: i,
			}

			if err := s.r.Create(ctx, lesson); err != nil {
				return nil, fmt.Errorf("failed to create lesson %d for module %s: %w", i, module.Title, err)
			}
			seededLessons = append(seededLessons, lesson)
		}
	}

	return seededLessons, nil
}
