package domain

import (
	"time"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/shared"
	sub "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/submission/domain"
)

type AssessmentType int
type AssessmentSubType int

const (
	Assignment AssessmentType = iota
	Exam
)

const (
	Exercise AssessmentSubType = iota
	Homework
	Quiz
	AssessmentExam
	MidtermExam
	PracticalExam
	FinalExam
)

type Assessment struct {
	shared.Base

	Title string
	Type AssessmentType
	SubType AssessmentSubType
	DueDate time.Time
	Submissions []sub.Submission
}