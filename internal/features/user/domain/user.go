package domain

import (
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/shared"
	"github.com/google/uuid"
)

type UserMetadata struct {
	Address string `json:"address"`
	BloodType string `json:"blood_type"`
}

type User struct {
	shared.Base

	Email string 
	PasswordHash string
	FirstName string
	LastName string
	ExternalID string // NISN for schools, NIM for University
	Metadata *UserMetadata
	GuardianID *uuid.UUID

	IsActive bool
}

func (u *User) IsChildOf(possibleGuardian User) bool {
	if u.GuardianID == nil {
		return false
	}
	return *u.GuardianID == possibleGuardian.Base.ID
}