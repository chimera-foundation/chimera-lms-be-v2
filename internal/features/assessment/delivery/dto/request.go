package dto

// GetStudentAssessmentsRequest represents the request parameters for fetching student assessments
type GetStudentAssessmentsRequest struct {
	Type      string `json:"type"`       // "assignment" or "exam"
	StartDate string `json:"start_date"` // RFC3339 format
	EndDate   string `json:"end_date"`   // RFC3339 format
	Limit     int    `json:"limit"`
	Offset    int    `json:"offset"`
}
