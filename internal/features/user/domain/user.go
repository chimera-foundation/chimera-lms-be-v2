package domain

import (
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/shared"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserMetadata struct {
	Address   string `json:"address"`
	BloodType string `json:"blood_type"`
}

type User struct {
	shared.Base

	OrganizationID uuid.UUID

	Email        string
	PasswordHash string
	FirstName    string
	LastName     string
	Metadata     *UserMetadata // TODO: add to migration
	GuardianID   *uuid.UUID // TODO: add to migration
	Roles        []Role

	IsSuperuser bool
}

func NewUser(email, firstName, lastName string, orgID uuid.UUID, roles []Role) *User {
	return &User{
		Email:          email,
		FirstName:      firstName,
		LastName:       lastName,
		OrganizationID: orgID,
		Roles: 			roles,
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

func (u *User) HasAnyRole(roleNames ...string) bool {
    for _, r := range u.Roles {
        for _, name := range roleNames {
            if r.Name == name {
                return true
            }
        }
    }
    return false
}

func (u *User) RolesStr() []string {
    if len(u.Roles) == 0 {
        return nil
    }

    roles := make([]string, len(u.Roles))
    for i := range u.Roles {
        roles[i] = u.Roles[i].Name
    }

    return roles
}