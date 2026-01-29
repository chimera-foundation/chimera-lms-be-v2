package domain

import (
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/shared"
)

type Role struct {
	shared.Base

	Name string
	Permissions map[string][]string `json:"permissions"`
}

func (r *Role) HasPermission(resource, action string) bool {
    actions, ok := r.Permissions[resource]
    if !ok {
        return false
    }
    for _, a := range actions {
        if a == action {
            return true
        }
    }
    return false
}