package repository

import (
	"context"

	"crowdreview/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AchievementRepository manages gamification data.
type AchievementRepository interface {
	List(ctx context.Context) ([]models.Achievement, error)
	Grant(ctx context.Context, userID, achievementID uuid.UUID) error
}

type GormAchievementRepository struct {
	db *gorm.DB
}

func (r *GormAchievementRepository) List(ctx context.Context) ([]models.Achievement, error) {
	var achievements []models.Achievement
	if err := r.db.WithContext(ctx).Find(&achievements).Error; err != nil {
		return nil, err
	}
	return achievements, nil
}

func (r *GormAchievementRepository) Grant(ctx context.Context, userID, achievementID uuid.UUID) error {
	ua := models.UserAchievement{
		UserID:        userID,
		AchievementID: achievementID,
	}
	return r.db.WithContext(ctx).Create(&ua).Error
}
