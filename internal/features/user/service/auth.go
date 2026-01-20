package service

import (
	"context"
	"errors"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/user/domain"
	"github.com/google/uuid"
)

type Auth interface {
    Register(ctx context.Context, email, password, firstName, lastName string, orgID uuid.UUID) (*domain.User, error)
    Login(ctx context.Context, email, password string) (string, error)
}

type authService struct {
    repo domain.UserRepository
}

func NewAuthService(r domain.UserRepository) Auth {
	return &authService{repo: r}
}

func (s *authService) Register(ctx context.Context, email, password, firstName, lastName string, orgID uuid.UUID) (*domain.User, error) {
    existing, _ := s.repo.GetByEmail(ctx, email)
    if existing != nil {
        return nil, errors.New("user with this email already exists")
    }

    user := &domain.User{
        Email: email,
        FirstName: firstName,
        LastName: lastName,
        OrganizationID: orgID,
        IsActive: true,
    }

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

    // 3. Logic for generating a JWT would go here.
    // Usually, you would inject a 'TokenProvider' interface into this struct.
    token := "sample-jwt-token" 

    return token, nil
}