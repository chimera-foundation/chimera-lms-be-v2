package domain

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID string
	Email string
	Username string
	PasswordHash string
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}