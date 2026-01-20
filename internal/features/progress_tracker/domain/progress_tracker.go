package domain

import (
	"time"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/shared"
)

type ProgressTracker struct {
	shared.Base

	IsCompleted bool
	UpdatedAt time.Time
}