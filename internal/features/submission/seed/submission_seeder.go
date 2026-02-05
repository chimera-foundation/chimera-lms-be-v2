package seed

import (
	"context"
	"fmt"
	"time"

	assessmentDomain "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/assessment/domain"
	enrollmentDomain "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/enrollment/domain"
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/submission/domain"
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/shared"
	"github.com/google/uuid"
)

type SubmissionSeeder struct {
	r domain.SubmissionRepository
}

func NewSubmissionSeeder(r domain.SubmissionRepository) *SubmissionSeeder {
	return &SubmissionSeeder{
		r: r,
	}
}

// SeedSubmissions creates submissions with varied statuses for each enrolled student:
// - done: submitted + graded
// - submitted: awaiting grade
// - pending: no submission (future due date) - we skip these
// - overdue: no submission (past due date) - we skip these
func (s *SubmissionSeeder) SeedSubmissions(
	ctx context.Context,
	enrollments []*enrollmentDomain.Enrollment,
	assessments []*assessmentDomain.Assessment,
) ([]*domain.Submission, error) {
	var seededSubmissions []*domain.Submission
	now := time.Now()

	// Build a map of courseID -> assessments
	assessmentsByCourse := make(map[uuid.UUID][]*assessmentDomain.Assessment)
	for _, a := range assessments {
		assessmentsByCourse[a.CourseID] = append(assessmentsByCourse[a.CourseID], a)
	}

	for _, enrollment := range enrollments {
		courseAssessments := assessmentsByCourse[enrollment.CourseID]
		if len(courseAssessments) == 0 {
			continue
		}

		for i, assessment := range courseAssessments {
			// Create varied submission statuses
			// Pattern: index 0 = done, index 1 = submitted, index 2+ = skip (pending/overdue)
			var submission *domain.Submission

			switch i % 5 {
			case 0: // done - submitted and graded (for past due assignments)
				if assessment.DueDate.Before(now) {
					submittedAt := assessment.DueDate.Add(-24 * time.Hour) // submitted 1 day before due
					submission = &domain.Submission{
						Base: shared.Base{
							ID:        uuid.New(),
							CreatedAt: now,
							UpdatedAt: now,
						},
						AssessmentID: assessment.ID,
						UserID:       enrollment.UserID,
						EnrollmentID: enrollment.ID,
						FinalScore:   85.5 + float32(i)*2, // varied scores
						SubmittedAt:  submittedAt,
					}
				}
			case 1: // submitted but not graded (awaiting grade)
				if assessment.DueDate.Before(now) {
					submittedAt := assessment.DueDate.Add(-12 * time.Hour)
					submission = &domain.Submission{
						Base: shared.Base{
							ID:        uuid.New(),
							CreatedAt: now,
							UpdatedAt: now,
						},
						AssessmentID: assessment.ID,
						UserID:       enrollment.UserID,
						EnrollmentID: enrollment.ID,
						FinalScore:   0, // not graded yet
						SubmittedAt:  submittedAt,
					}
				}
			case 2: // overdue - no submission, past due date
				// Skip - no submission needed
			case 3: // done - another graded one for exams
				if assessment.DueDate.Before(now) {
					submittedAt := assessment.DueDate.Add(-2 * time.Hour)
					submission = &domain.Submission{
						Base: shared.Base{
							ID:        uuid.New(),
							CreatedAt: now,
							UpdatedAt: now,
						},
						AssessmentID: assessment.ID,
						UserID:       enrollment.UserID,
						EnrollmentID: enrollment.ID,
						FinalScore:   92.0 + float32(i),
						SubmittedAt:  submittedAt,
					}
				}
			case 4: // pending - no submission, future due date
				// Skip - no submission needed
			}

			if submission != nil {
				if err := s.r.Create(ctx, submission); err != nil {
					return nil, fmt.Errorf("failed to create submission for assessment '%s': %w", assessment.Title, err)
				}
				seededSubmissions = append(seededSubmissions, submission)
			}
		}
	}

	return seededSubmissions, nil
}
