package adapters

import (
	"context"
    "github.com/uptrace/bun"
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/modules/user/domain"
    "github.com/google/uuid"
)

type UserSQL struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID       uuid.UUID  `bun:",pk,autoincrement"`
    Username string `bun:",notnull"`
	Email    string `bun:",unique,notnull"`
	Password string `bun:",notnull"`
}

type BunUserRepo struct {
	db *bun.DB
}

func NewBunUserRepo(db *bun.DB) *BunUserRepo {
    _, _ = db.NewCreateTable().Model((*UserSQL)(nil)).IfNotExists().Exec(context.Background())
    if err := db.Ping(); err != nil {
        panic("CANNOT CONNECT TO POSTGRES: " + err.Error())
    }
	return &BunUserRepo{db: db}
}

func (r *BunUserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
    u := &domain.User{}
    query := "SELECT id, email, username, password_hash FROM users WHERE email = $1"
    err := r.db.QueryRowContext(ctx, query, email).Scan(&u.ID, &u.Email, &u.Username, &u.PasswordHash)
    return u, err
}

func (r *BunUserRepo) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
    u := &domain.User{}
    query := "SELECT id, email, username, password_hash FROM users WHERE username = $1"
    err := r.db.QueryRowContext(ctx, query, username).Scan(&u.ID, &u.Email, &u.Username, &u.PasswordHash)
    return u, err
}