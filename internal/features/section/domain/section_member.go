package domain

import (
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/shared"
	"github.com/google/uuid"
)


type SectionRoleType string

const (
	Student SectionRoleType = "student"
	Teacher SectionRoleType = "teacher"
	Assistant SectionRoleType = "assistant"
	Monitor SectionRoleType = "monitor"
)

type SectionMember struct {
	shared.Base

	SectionID uuid.UUID
	UserID uuid.UUID
	Type SectionRoleType 
}