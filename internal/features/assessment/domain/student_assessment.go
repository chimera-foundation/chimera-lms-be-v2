package domain

import (
	"time"

	"github.com/google/uuid"
)

// SubmissionStatus represents the status of a student's assessment submission
type SubmissionStatus string

const (
	StatusPending   SubmissionStatus = "pending"
	StatusSubmitted SubmissionStatus = "submitted"
	StatusDone      SubmissionStatus = "done"
	StatusOverdue   SubmissionStatus = "overdue"
)

// StudentAssessmentFilter contains the filter criteria for fetching student assessments
type StudentAssessmentFilter struct {
	Type      *AssessmentType // Optional: filter by assessment type
	StartDate *time.Time      // Optional: filter by due date >= start
	EndDate   *time.Time      // Optional: filter by due date <= end
	Limit     int
	Offset    int
}

// StudentAssessmentItem represents an assessment with its submission status for a student
type StudentAssessmentItem struct {
	ID            uuid.UUID
	Subject       string
	Title         string
	AttachmentURL *string
	Status        SubmissionStatus
	Type          AssessmentType
	SubType       AssessmentSubType
	DueDate       time.Time
}

// StudentAssessmentSummary contains the counts of assessments by status
type StudentAssessmentSummary struct {
	Pending   int
	Submitted int
	Done      int
	Overdue   int
}
