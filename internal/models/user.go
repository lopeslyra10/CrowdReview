package models

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// User represents an end-user with optional gamification data.
type User struct {
	Base
	Email            string            `gorm:"uniqueIndex:idx_users_email,where:deleted_at IS NULL;not null"`
	Username         string            `gorm:"uniqueIndex:idx_users_username,where:deleted_at IS NULL;not null"`
	PasswordHash     string            `gorm:"not null"`
	Role             string            `gorm:"type:varchar(20);index;default:'user'"` // user or admin
	GamificationScore int              `gorm:"default:0"`
	ProfileMeta      datatypes.JSONMap `gorm:"type:jsonb;default:'{}'::jsonb"`
	Achievements     []UserAchievement
	Reviews          []Review
}

// AdminUser stores admin-only metadata while sharing credentials with User.
type AdminUser struct {
	Base
	UserID      uuid.UUID        `gorm:"type:uuid;uniqueIndex"`
	User        User             `gorm:"constraint:OnDelete:CASCADE"`
	Permissions datatypes.JSONMap `gorm:"type:jsonb;default:'{}'::jsonb"`
}
