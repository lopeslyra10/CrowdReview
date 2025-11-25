package repository

import (
	"context"

	"crowdreview/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ValidationRepository persists validation results.
type ValidationRepository interface {
	SaveResult(ctx context.Context, result *models.ReviewValidationResult) error
	MarkReview(ctx context.Context, reviewID uuid.UUID, resultID uuid.UUID, status string, suspicious bool) error
}

type GormValidationRepository struct {
	db *gorm.DB
}

func (r *GormValidationRepository) SaveResult(ctx context.Context, result *models.ReviewValidationResult) error {
	return r.db.WithContext(ctx).Create(result).Error
}

func (r *GormValidationRepository) MarkReview(ctx context.Context, reviewID uuid.UUID, resultID uuid.UUID, status string, suspicious bool) error {
	return r.db.WithContext(ctx).
		Model(&models.Review{}).
		Where("id = ?", reviewID).
		Updates(map[string]interface{}{
			"validation_result_id": resultID,
			"status":               status,
			"suspicious":           suspicious,
		}).Error
}
