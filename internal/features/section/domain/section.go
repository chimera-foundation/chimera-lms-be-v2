package domain

import (
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/shared"
)

type Section struct {
	shared.Base

	Name string
	RoomCode string
	Capacity int
}