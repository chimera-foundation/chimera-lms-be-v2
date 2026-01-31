package seed

import (
	"context"
	"fmt"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/user/domain"
	"github.com/google/uuid"
)

type UserSeeder struct {
	ur domain.UserRepository
	rr domain.RoleRepository
}

func NewUserSeeder(ur domain.UserRepository, rr domain.RoleRepository) *UserSeeder {
	return &UserSeeder{
		ur: ur,
		rr: rr,
	}
}

func (s *UserSeeder) SeedUsers(ctx context.Context, orgID uuid.UUID) (map[string]*domain.User, error) {
	// 1. Fetch Roles
	adminRole, err := s.rr.GetByName(ctx, "admin")
	if err != nil {
		return nil, fmt.Errorf("failed to get admin role: %w", err)
	}
	teacherRole, err := s.rr.GetByName(ctx, "teacher")
	if err != nil {
		return nil, fmt.Errorf("failed to get teacher role: %w", err)
	}
	studentRole, err := s.rr.GetByName(ctx, "student")
	if err != nil {
		return nil, fmt.Errorf("failed to get student role: %w", err)
	}

	// 2. Define Users
	usersToSeed := []struct {
		Email     string
		FirstName string
		LastName  string
		Role      *domain.Role
		Grade     string // Optional, for identification
	}{
		{
			Email:     "admin@candletree.com",
			FirstName: "Admin",
			LastName:  "User",
			Role:      adminRole,
		},
		{
			Email:     "teacher@candletree.com",
			FirstName: "Teacher",
			LastName:  "User",
			Role:      teacherRole,
		},
		{
			Email:     "student10@candletree.com",
			FirstName: "Student",
			LastName:  "Kelas 10",
			Role:      studentRole,
			Grade:     "10",
		},
		{
			Email:     "student11@candletree.com",
			FirstName: "Student",
			LastName:  "Kelas 11",
			Role:      studentRole,
			Grade:     "11",
		},
	}

	seededUsers := make(map[string]*domain.User)

	for _, uVal := range usersToSeed {
		// Check if user exists
		existing, err := s.ur.GetByEmail(ctx, uVal.Email)
		if err != nil {
			return nil, fmt.Errorf("failed to check existing user %s: %w", uVal.Email, err)
		}

		if existing != nil {
			seededUsers[uVal.Email] = existing
			continue
		}

		// Create new user
		user := domain.NewUser(
			uVal.Email,
			uVal.FirstName,
			uVal.LastName,
			orgID,
			[]domain.Role{*uVal.Role},
		)

		if err := user.SetPassword("password"); err != nil {
			return nil, fmt.Errorf("failed to set password for %s: %w", uVal.Email, err)
		}

		if err := s.ur.Create(ctx, user); err != nil {
			return nil, fmt.Errorf("failed to create user %s: %w", uVal.Email, err)
		}

		seededUsers[uVal.Email] = user
	}

	return seededUsers, nil
}
