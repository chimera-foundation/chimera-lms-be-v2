package domain

import (
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/shared"
)

type Program struct {
	shared.Base

	Name string
	Description string
}