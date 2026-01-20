package domain

import (
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/shared"
	"github.com/google/uuid"
)


type SectionRoleType int

const (
	Student SectionRoleType = iota
	Teacher
	Assistant
	Monitor
)

type SectionMember struct {
	shared.Base

	SectionID uuid.UUID
	UserID uuid.UUID
	Type SectionRoleType 
}