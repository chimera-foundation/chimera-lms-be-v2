package domain

import (
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/shared"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserMetadata struct {
	Address string `json:"address"`
	BloodType string `json:"blood_type"`
}

type User struct {
	shared.Base

	OrganizationID uuid.UUID

	Email string 
	PasswordHash string
	FirstName string
	LastName string
	Metadata *UserMetadata
	GuardianID *uuid.UUID
	Roles []Role

	IsActive bool
	IsSuperuser bool
}

func NewUser(email, firstName, lastName string, orgID uuid.UUID) *User {
    return &User{
        Email:          email,
        FirstName:      firstName,
        LastName:       lastName,
        OrganizationID: orgID,
        IsActive:       true, 
        IsSuperuser:    false, 
        Metadata:       &UserMetadata{}, 
    }
}

func (u *User) IsChildOf(possibleGuardian User) bool {
	if u.GuardianID == nil {
		return false
	}
	return *u.GuardianID == possibleGuardian.Base.ID
}

func (u *User) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (u *User) SetPassword(password string) error {
    bytes, err := u.HashPassword(password)
    if err != nil {
        return err
    }
    u.PasswordHash = string(bytes)
    return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}