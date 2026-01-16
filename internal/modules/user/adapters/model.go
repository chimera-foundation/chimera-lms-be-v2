package adapters

import (
    "github.com/uptrace/bun"
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/modules/user/domain"
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/shared"
    "github.com/google/uuid"
)

type UserSQL struct {
	bun.BaseModel `bun:"table:users,alias:u"`
	
	ID       uuid.UUID  `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
    Username string `bun:"username,notnull"`
	Email    string `bun:"email,unique,notnull"`
	PasswordHash string `bun:"password_hash,notnull"`
	
	shared.AuditModel
}

func (u *UserSQL) ToDomain() *domain.User {
	return &domain.User{
		ID:           u.ID.String(),
		Email:        u.Email,
		Username:     u.Username,
		PasswordHash: u.PasswordHash,
	}
}