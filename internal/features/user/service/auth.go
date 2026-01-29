package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/user/domain"
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/shared/auth"
	"github.com/google/uuid"
)

type authService struct {
    userRepo domain.UserRepository
    roleRepo domain.RoleRepository
    tokenProvider auth.TokenProvider
}

func NewAuthService(ur domain.UserRepository, rr domain.RoleRepository, tp auth.TokenProvider) Auth {
    return &authService{
        userRepo:  ur,
        roleRepo:  rr,
        tokenProvider: tp,
    }
}

func (s *authService) RegisterStudent(ctx context.Context, email, password, firstName, lastName string, orgID uuid.UUID) (*domain.User, error) {
    studentRole, err := s.roleRepo.GetByName(ctx, "student")
    if err != nil {
        return nil, fmt.Errorf("system error: roles not configured")
    }
    if studentRole == nil {
        return nil, errors.New("registration failed: student role does not exist")
    }

    user := domain.NewUser(
        email, 
        firstName, 
        lastName, 
        orgID, 
        []domain.Role{*studentRole},
    )

    if err := user.SetPassword(password); err != nil {
        return nil, err
    }

    if err := s.userRepo.Create(ctx, user); err != nil {
        return nil, errors.New("user with this email already exists")
    }

    return user, nil
}

func (s *authService) Login(ctx context.Context, email, password string) (string, error) {
    user, err := s.userRepo.GetByEmail(ctx, email)
    if user==nil || err != nil || !user.CheckPassword(password) {
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

    user, err := s.userRepo.GetByID(ctx, user_id)
    if err != nil {
        return nil, errors.New("user does not exists")
    }
    
    return user, nil
}