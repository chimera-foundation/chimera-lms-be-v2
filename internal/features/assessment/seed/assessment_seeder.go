package seed

import (
	"context"
	"fmt"
	"time"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/assessment/domain"
	contentDomain "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/content/domain"
	courseDomain "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/course/domain"
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/shared"
	"github.com/google/uuid"
)

type AssessmentSeeder struct {
	assessmentRepo domain.AssessmentRepo
	contentRepo    contentDomain.ContentRepository
}

func NewAssessmentSeeder(assessmentRepo domain.AssessmentRepo, contentRepo contentDomain.ContentRepository) *AssessmentSeeder {
	return &AssessmentSeeder{
		assessmentRepo: assessmentRepo,
		contentRepo:    contentRepo,
	}
}

// SeedAssessments creates assessments for each course:
// - 3 assignments (exercise, homework, quiz) with varied due dates
// - 2 exams (midterm, final)
// Each assessment also gets a content record with attachment URL
func (s *AssessmentSeeder) SeedAssessments(
	ctx context.Context,
	courses []*courseDomain.Course,
) ([]*domain.Assessment, error) {
	var seededAssessments []*domain.Assessment
	now := time.Now()

	for _, course := range courses {
		// Define assessments for this course
		assessmentsToCreate := []struct {
			title   string
			aType   domain.AssessmentType
			subType domain.AssessmentSubType
			dueDate time.Time
		}{
			// Assignments - varied due dates for different statuses
			{
				title:   fmt.Sprintf("%s - Exercise 1", course.Title),
				aType:   domain.Assignment,
				subType: domain.Exercise,
				dueDate: now.AddDate(0, 0, -14), // 2 weeks ago (for overdue)
			},
			{
				title:   fmt.Sprintf("%s - Homework 1", course.Title),
				aType:   domain.Assignment,
				subType: domain.Homework,
				dueDate: now.AddDate(0, 0, -7), // 1 week ago (for done/submitted)
			},
			{
				title:   fmt.Sprintf("%s - Quiz 1", course.Title),
				aType:   domain.Assignment,
				subType: domain.Quiz,
				dueDate: now.AddDate(0, 0, 7), // 1 week from now (for pending)
			},
			// Exams
			{
				title:   fmt.Sprintf("%s - Midterm Exam", course.Title),
				aType:   domain.Exam,
				subType: domain.MidtermExam,
				dueDate: now.AddDate(0, 0, -3), // 3 days ago (for submitted)
			},
			{
				title:   fmt.Sprintf("%s - Final Exam", course.Title),
				aType:   domain.Exam,
				subType: domain.FinalExam,
				dueDate: now.AddDate(0, 0, 30), // 1 month from now (for pending)
			},
		}

		for i, a := range assessmentsToCreate {
			assessment := &domain.Assessment{
				Base: shared.Base{
					ID:        uuid.New(),
					CreatedAt: now,
					UpdatedAt: now,
				},
				OrganizationID: course.OrganizationID,
				CourseID:       course.ID,
				Title:          a.title,
				Type:           a.aType,
				SubType:        a.subType,
				DueDate:        a.dueDate,
			}

			if err := s.assessmentRepo.Create(ctx, assessment); err != nil {
				return nil, fmt.Errorf("failed to create assessment '%s': %w", a.title, err)
			}

			// Create content record with attachment URL for each assessment
			content := &contentDomain.Content{
				Base: shared.Base{
					ID:        uuid.New(),
					CreatedAt: now,
					UpdatedAt: now,
				},
				AssessmentID: assessment.ID,
				Type:         contentDomain.Document,
				Data: &contentDomain.ContentData{
					URL:         fmt.Sprintf("https://storage.lms.example.com/attachments/%s/%d.pdf", course.ID.String(), i+1),
					Title:       fmt.Sprintf("%s - Instructions", a.title),
					Description: fmt.Sprintf("Attachment for %s", a.title),
				},
			}

			if err := s.contentRepo.Create(ctx, content); err != nil {
				return nil, fmt.Errorf("failed to create content for assessment '%s': %w", a.title, err)
			}

			seededAssessments = append(seededAssessments, assessment)
		}
	}

	return seededAssessments, nil
}
