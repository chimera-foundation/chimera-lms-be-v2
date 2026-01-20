package domain

import (
	"time"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/shared"
)

type Submission struct {
	shared.Base

	FinalScore float32
	SubmittedAt time.Time
}