package dto

type RegisterRequest struct {
	Email          string    `json:"email"`
	Password       string    `json:"password"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	OrganizationID string `json:"organization_id"`
}

type RegisterResponse struct {
	Email string
	FirstName string
	LastName string
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type MeResponse struct {
	Email          string    `json:"email"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	Roles          []string  `json:"roles"`
}