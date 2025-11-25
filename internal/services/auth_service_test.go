package services

import (
	"context"
	"testing"
	"time"

	"crowdreview/config"
	"crowdreview/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type mockUserRepo struct {
	users map[string]*models.User
}

func (m *mockUserRepo) Create(ctx context.Context, user *models.User) error {
	if m.users == nil {
		m.users = make(map[string]*models.User)
	}
	m.users[user.Email] = user
	return nil
}
func (m *mockUserRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	if u, ok := m.users[email]; ok {
		return u, nil
	}
	return nil, gorm.ErrRecordNotFound
}
func (m *mockUserRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	for _, u := range m.users {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func TestAuthServiceRegisterAndLogin(t *testing.T) {
	repo := &mockUserRepo{users: map[string]*models.User{}}
	cfg := config.Config{
		JWTSecret:     "secret",
		RefreshSecret: "refresh",
		TokenTTL:      time.Minute,
		RefreshTTL:    time.Hour,
	}
	service := &DefaultAuthService{Users: repo, Config: cfg}

	user, access, refresh, err := service.Register(context.Background(), "a@b.com", "testuser", "password123")
	require.NoError(t, err)
	require.NotEmpty(t, access)
	require.NotEmpty(t, refresh)
	require.Equal(t, "a@b.com", user.Email)

	_, access2, _, err := service.Login(context.Background(), "a@b.com", "password123")
	require.NoError(t, err)
	require.NotEmpty(t, access2)
}
