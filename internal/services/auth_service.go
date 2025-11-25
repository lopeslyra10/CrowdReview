package services

import (
	"context"
	"errors"
	"time"

	"crowdreview/config"
	"crowdreview/internal/models"
	"crowdreview/internal/repository"
	"crowdreview/pkg/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AuthService exposes auth flows.
type AuthService interface {
	Register(ctx context.Context, email, username, password string) (*models.User, string, string, error)
	Login(ctx context.Context, email, password string) (*models.User, string, string, error)
	Refresh(ctx context.Context, userID uuid.UUID) (string, string, error)
	ValidateRefreshToken(token string) (uuid.UUID, error)
}

type DefaultAuthService struct {
	Users  repository.UserRepository
	Config config.Config
}

func (s *DefaultAuthService) Register(ctx context.Context, email, username, password string) (*models.User, string, string, error) {
	hash, err := utils.HashPassword(password)
	if err != nil {
		return nil, "", "", err
	}
	user := &models.User{
		Email:        email,
		Username:     username,
		PasswordHash: hash,
		Role:         "user",
	}
	if err := s.Users.Create(ctx, user); err != nil {
		return nil, "", "", err
	}
	access, refresh, err := utils.GenerateTokens(user.ID, user.Role, s.Config)
	return user, access, refresh, err
}

func (s *DefaultAuthService) Login(ctx context.Context, email, password string) (*models.User, string, string, error) {
	user, err := s.Users.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", "", errors.New("invalid credentials")
		}
		return nil, "", "", err
	}
	if err := utils.VerifyPassword(user.PasswordHash, password); err != nil {
		return nil, "", "", errors.New("invalid credentials")
	}
	access, refresh, err := utils.GenerateTokens(user.ID, user.Role, s.Config)
	return user, access, refresh, err
}

func (s *DefaultAuthService) Refresh(ctx context.Context, userID uuid.UUID) (string, string, error) {
	user, err := s.Users.GetByID(ctx, userID)
	if err != nil {
		return "", "", err
	}
	access, refresh, err := utils.GenerateTokens(user.ID, user.Role, s.Config)
	return access, refresh, err
}

// ValidateRefreshToken ensures refresh token is valid; kept here for convenience.
func (s *DefaultAuthService) ValidateRefreshToken(token string) (uuid.UUID, error) {
	claims, err := utils.ParseToken(token, s.Config.RefreshSecret)
	if err != nil {
		return uuid.Nil, err
	}
	sub, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, err
	}
	if time.Unix(claims.ExpiresAt.Unix(), 0).Before(time.Now()) {
		return uuid.Nil, errors.New("refresh token expired")
	}
	return sub, nil
}
