package domain

import (
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/shared"
)

type Role struct {
	shared.Base

	Name string
}