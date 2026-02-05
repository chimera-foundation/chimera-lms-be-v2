package domain

import (
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