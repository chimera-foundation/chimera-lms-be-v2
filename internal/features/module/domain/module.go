package domain

import (
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/shared"
)

type Module struct {
	shared.Base

	Title string
	OrderIndex int
}