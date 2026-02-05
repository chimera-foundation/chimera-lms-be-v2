package dto

import (
	"time"

	"github.com/google/uuid"
)

// AssessmentSummary contains counts of assessments by status
type AssessmentSummary struct {
	Pending   int `json:"pending"`
	Submitted int `json:"submitted"`
	Done      int `json:"done"`
	Overdue   int `json:"overdue"`
}

// AssessmentItem represents a single assessment in the student's list
type AssessmentItem struct {
	ID            uuid.UUID `json:"id"`
	Subject       string    `json:"subject"`
	Title         string    `json:"title"`
	AttachmentURL *string   `json:"attachment_url,omitempty"`
	Status        string    `json:"status"`
	Type          string    `json:"type"`
	SubType       string    `json:"sub_type"`
	DueDate       time.Time `json:"due_date"`
}

// StudentAssessmentsResponse is the response structure for the student assessments endpoint
type StudentAssessmentsResponse struct {
	Summary     AssessmentSummary `json:"summary"`
	Assessments []AssessmentItem  `json:"assessments"`
}
