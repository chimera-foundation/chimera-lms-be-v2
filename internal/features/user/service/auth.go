package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/user/domain"
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/shared/auth"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type authService struct {
	userRepo      domain.UserRepository
	roleRepo      domain.RoleRepository
	tokenProvider auth.TokenProvider
	log           *logrus.Logger
}

func NewAuthService(ur domain.UserRepository, rr domain.RoleRepository, tp auth.TokenProvider, log *logrus.Logger) Auth {
	return &authService{
		userRepo:      ur,
		roleRepo:      rr,
		tokenProvider: tp,
		log:           log,
	}
}

func (s *authService) register(ctx context.Context, email, password, firstName, lastName string, orgID uuid.UUID, roleName string) (*domain.User, error) {
	role, err := s.roleRepo.GetByName(ctx, roleName)
	if err != nil {
		s.log.WithError(err).WithField("role_name", roleName).Error("failed to get role during registration")
		return nil, fmt.Errorf("system error: %w", err)
	}
	if role == nil {
		s.log.WithField("role_name", roleName).Warn("role does not exist during registration")
		return nil, fmt.Errorf("registration failed: role '%s' does not exist", roleName)
	}

	user := domain.NewUser(email, firstName, lastName, orgID, []domain.Role{*role})

	if err := user.SetPassword(password); err != nil {
		return nil, err
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		s.log.WithError(err).WithField("email", email).Warn("user with this email already exists")
		return nil, errors.New("user with this email already exists")
	}

	s.log.WithFields(logrus.Fields{"user_id": user.ID, "email": email, "role": roleName}).Info("user registered successfully")
	return user, nil
}

func (s *authService) RegisterStudent(ctx context.Context, email, password, firstName, lastName string, orgID uuid.UUID) (*domain.User, error) {
	return s.register(ctx, email, password, firstName, lastName, orgID, "student")
}

func (s *authService) RegisterTeacher(ctx context.Context, email, password, firstName, lastName string, orgID uuid.UUID) (*domain.User, error) {
	return s.register(ctx, email, password, firstName, lastName, orgID, "teacher")
}

func (s *authService) RegisterAdmin(ctx context.Context, email, password, firstName, lastName string, orgID uuid.UUID) (*domain.User, error) {
	return s.register(ctx, email, password, firstName, lastName, orgID, "admin")
}

func (s *authService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if user == nil || err != nil || !user.CheckPassword(password) {
		s.log.WithField("email", email).Warn("login failed: invalid email or password")
		return "", errors.New("invalid email or password")
	}

	token, err := s.tokenProvider.GenerateToken(user.ID, user.OrganizationID)
	if err != nil {
		s.log.WithError(err).WithField("user_id", user.ID).Error("failed to generate access token")
		return "", errors.New("failed to generate access token")
	}

	s.log.WithField("user_id", user.ID).Info("user logged in successfully")
	return token, nil
}

func (s *authService) Logout(ctx context.Context, token string) error {
	err := s.tokenProvider.BlacklistToken(ctx, token, 15*time.Minute)
	if err != nil {
		s.log.WithError(err).Error("failed to blacklist token during logout")
		return errors.New("could not invalidate session")
	}

	s.log.Info("user logged out successfully")
	return nil
}

func (s *authService) Me(ctx context.Context) (*domain.User, error) {
	user_id, ok := auth.GetUserID(ctx)
	if !ok {
		s.log.Warn("failed getting user id from context")
		return nil, errors.New("failed getting user id")
	}

	user, err := s.userRepo.GetByID(ctx, user_id)
	if err != nil {
		s.log.WithError(err).WithField("user_id", user_id).Error("failed to get user by id")
		return nil, errors.New("user does not exists")
	}

	return user, nil
}
