package domain

import (
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/shared"
	con "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/content/domain"
)

type Lesson struct {
	shared.Base

	Title string
	OrderIndex int
	Contents []con.Content
}