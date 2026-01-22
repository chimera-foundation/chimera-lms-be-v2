package dto

import "github.com/google/uuid"

type RegisterRequest struct {
	Email          string    `json:"email"`
	Password       string    `json:"password"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	OrganizationID uuid.UUID `json:"organization_id"`
}

type LoginRequest struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }