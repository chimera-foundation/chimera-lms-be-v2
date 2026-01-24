package service

import (
	"context"
	"errors"
	"time"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/user/domain"
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/shared/auth"
	"github.com/google/uuid"
)

type authService struct {
    repo domain.UserRepository
    tokenProvider auth.TokenProvider
}

func NewAuthService(r domain.UserRepository, tp auth.TokenProvider) Auth {
	return &authService{
        repo: r,
        tokenProvider: tp,
    }
}

func (s *authService) Register(ctx context.Context, email, password, firstName, lastName string, orgID uuid.UUID) (*domain.User, error) {
    existing, _ := s.repo.GetByEmail(ctx, email)
    if existing != nil {
        return nil, errors.New("user with this email already exists")
    }

    user := domain.NewUser(
        email,
        firstName,
        lastName,
        orgID,
    )

    if err := user.SetPassword(password); err != nil {
        return nil, err
    }

    if err := s.repo.Create(ctx, user); err != nil {
        return nil, err
    }

    return user, nil
}

func (s *authService) Login(ctx context.Context, email, password string) (string, error) {
    user, err := s.repo.GetByEmail(ctx, email)
    if err != nil {
        return "", errors.New("invalid email or password")
    }

    if !user.CheckPassword(password) {
        return "", errors.New("invalid email or password")
    }

    token, err := s.tokenProvider.GenerateToken(user.ID, user.OrganizationID)
    if err != nil{
        return "", errors.New("failed to generate access token")
    }

    return token, nil
}

func (s *authService) Logout(ctx context.Context, token string) (error) {
    err := s.tokenProvider.BlacklistToken(ctx, token, 15 * time.Minute) 
    if err != nil {
        return errors.New("could not invalidate session")
    }

    return nil
}

func (s *authService) Me(ctx context.Context, token string) (*domain.User, error) {
    user_id, err := s.tokenProvider.ValidateToken(token)
    if err != nil {
        return nil, errors.New("invalid token")
    }

    user, err := s.repo.GetByID(ctx, user_id)
    if err != nil {
        return nil, errors.New("user does not exists")
    }
    
    return user, nil
}