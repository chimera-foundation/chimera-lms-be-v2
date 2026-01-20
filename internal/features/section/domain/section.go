package domain

import (
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/shared"
	u "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/user/domain"
)

type SectionRoleType int

const (
	Student SectionRoleType = iota
	Teacher
	Assistant
	Monitor
)

type Section struct {
	shared.Base

	Name string
	RoomCode string
	Capacity int

	Users []u.User
}