package domain

import (
	"github.com/google/uuid"
)

type SectionRoleType string

const (
	Student   SectionRoleType = "student"
	Teacher   SectionRoleType = "teacher"
	Assistant SectionRoleType = "assistant"
	Monitor   SectionRoleType = "monitor"
)

type SectionMember struct {
	SectionID uuid.UUID
	UserID    uuid.UUID
	RoleType  SectionRoleType
}
