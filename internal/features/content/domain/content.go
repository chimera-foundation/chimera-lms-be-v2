package domain

import (
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/shared"
	prog "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/progress_tracker/domain"
)

type ContentType int

const (
	Video ContentType = iota
	Document
	Quiz
)

type Content struct {
	shared.Base

	Type ContentType
	ProgressTrackers []prog.ProgressTracker
}