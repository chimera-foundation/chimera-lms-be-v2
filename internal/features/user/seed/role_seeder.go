package seed

import (
	"context"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/user/domain"
)

type RoleSeeder struct {
	rr domain.RoleRepository
}

func NewRoleSeeder(rr domain.RoleRepository) *RoleSeeder {
	return &RoleSeeder{
		rr: rr,
	}
}

func (s *RoleSeeder) SeedRoles(ctx context.Context) error {
	roles := []*domain.Role{
		{
			Name: "superadmin",
			Permissions: map[string][]string{
				"all": {"manage"}, 
			},
		},
		{
			Name: "admin",
			Permissions: map[string][]string{
				"user":         {"create", "read", "update", "delete", "impersonate"},
				"course":       {"read", "delete", "archive"},
				"organization": {"read", "update"},
				"report":       {"read", "export"},
				"billing":      {"read", "refund"},
			},
		},
		{
			Name: "teacher",
			Permissions: map[string][]string{
				"course": {"create", "read", "update", "delete"},
				"content": {"upload", "organize"},
				"student": {"grade", "view_progress"},
			},
		},
		{
			Name: "student",
			Permissions: map[string][]string{
				"course": {"read", "enroll"},
				"quiz":   {"take", "view_results"},
			},
		},
	}

	for _, role := range roles {
		existing, err := s.rr.GetByName(ctx, role.Name)
		if err != nil {
			return err
		}

		if existing == nil {
			if err := s.rr.Create(ctx, role); err != nil {
				return err
			}
		}
	}
	return nil
}