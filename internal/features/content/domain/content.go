package domain

import (
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/shared"
	"github.com/google/uuid"
)

type ContentType string

const (
	Video ContentType = "video"
	Document ContentType = "document"
	Quiz ContentType = "quiz"
)

type Content struct {
	shared.Base

	LessonID uuid.UUID
	AssessmentID uuid.UUID

	Type ContentType
}