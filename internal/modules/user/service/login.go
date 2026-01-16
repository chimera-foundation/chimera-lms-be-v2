package service

import (
	"context"
	"errors"
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/modules/user/ports"
	"github.com/golang-jwt/jwt/v5"
    "time"
)

type userService struct {
	repo ports.Repository
	jwtSecret []byte
}

func NewUserService(repo ports.Repository, secret string) ports.UseCase {
	return &userService{repo: repo, jwtSecret: []byte(secret)}
}

func (s *userService) LoginByEmail(ctx context.Context, email, password string) (string, error) {
    user, err := s.repo.GetByEmail(ctx, email)
    if err != nil || !user.CheckPassword(password) {
        return "", errors.New("invalid credentials")
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "sub": user.ID,
        "exp": time.Now().Add(time.Hour * 72).Unix(),
    })	

    return token.SignedString(s.jwtSecret)
}

func (s *userService) LoginByUsername(ctx context.Context, username, password string) (string, error) {
    user, err := s.repo.GetByUsername(ctx, username)
    if err != nil || !user.CheckPassword(password) {
        return "", errors.New("invalid credentials")
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "sub": user.ID,
        "exp": time.Now().Add(time.Hour * 72).Unix(),
    })	

    return token.SignedString(s.jwtSecret)
}