package models

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// Achievement describes gamification achievements.
type Achievement struct {
	Base
	Name        string            `gorm:"uniqueIndex:idx_achievements_name"`
	Description string
	Points      int               `gorm:"index"`
	Meta        datatypes.JSONMap `gorm:"type:jsonb;default:'{}'::jsonb"`
	UserAchievements []UserAchievement
}

// UserAchievement joins users and achievements.
type UserAchievement struct {
	Base
	UserID        uuid.UUID `gorm:"type:uuid;index"`
	AchievementID uuid.UUID `gorm:"type:uuid;index"`
	EarnedAt      int64     `gorm:"autoCreateTime"`
}
