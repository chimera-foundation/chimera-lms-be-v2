package adapters

import (
	"context"
    "github.com/uptrace/bun"
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/modules/user/domain"
)

type BunUserRepo struct {
	db *bun.DB
}

func NewBunUserRepo(db *bun.DB) *BunUserRepo {
    _, err := db.NewCreateTable().Model((*UserSQL)(nil)).IfNotExists().Exec(context.Background())
    if err != nil {
        panic("FAILED TO CREATE TABLE: " + err.Error())
    }
	return &BunUserRepo{db: db}
}

func (r *BunUserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
    var model UserSQL
	
	err := r.db.NewSelect().
		Model(&model).
		Where("email = ?", email).
		Scan(ctx)

	if err != nil {
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *BunUserRepo) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
    var model UserSQL
	
	err := r.db.NewSelect().
		Model(&model).
		Where("username = ?", username).
		Scan(ctx)

	if err != nil {
		return nil, err
	}
	return model.ToDomain(), nil
}