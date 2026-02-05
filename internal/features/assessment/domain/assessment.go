package domain

import (
	"errors"
	"time"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/shared"
	"github.com/google/uuid"
)

type AssessmentType string
type AssessmentSubType string

const (
	Assignment AssessmentType = "assignment"
	Exam AssessmentType = "exam"
)

const (
	Exercise AssessmentSubType = "exercise"
	Homework AssessmentSubType = "homework"
	Quiz AssessmentSubType = "quiz"
	AssessmentExam AssessmentSubType = "assessment_exam"
	MidtermExam AssessmentSubType = "midterm_exam"
	PracticalExam AssessmentSubType = "practical_exam"
	FinalExam AssessmentSubType = "final_exam"
)

type Assessment struct {
	shared.Base

	OrganizationID uuid.UUID
	CourseID uuid.UUID

	Title string
	Type AssessmentType
	SubType AssessmentSubType
	DueDate time.Time
}

func NewAssessment(
    orgID uuid.UUID,
    courseID uuid.UUID,
    title string,
    assessmentType string,
    subType string,
    dueDate time.Time,
) *Assessment {
    return &Assessment{
        OrganizationID: orgID,
        CourseID:       courseID,
        Title:          title,
        Type:           AssessmentType(assessmentType),
        SubType:        AssessmentSubType(subType),
        DueDate:        dueDate,
    }
}

func (a *Assessment) Validate() error {
    if a.Title == "" {
        return errors.New("assessment title is required")
    }
    if a.OrganizationID == uuid.Nil || a.CourseID == uuid.Nil {
        return errors.New("organization and course IDs must be valid")
    }
    return nil
}