package domain

import (
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/shared"
	"github.com/google/uuid"
)

type ContentType int

const (
	Video ContentType = iota
	Document
	Quiz
)

type Content struct {
	shared.Base

	LessonID uuid.UUID
	AssessmentID uuid.UUID

	Type ContentType
}