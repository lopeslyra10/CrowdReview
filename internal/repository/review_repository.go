package repository

import (
	"context"

	"crowdreview/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ReviewRepository stores reviews and aggregates.
type ReviewRepository interface {
	Create(ctx context.Context, review *models.Review) error
	ListByCompany(ctx context.Context, companyID uuid.UUID) ([]models.Review, error)
	ListSuspicious(ctx context.Context) ([]models.Review, error)
	Respond(ctx context.Context, id uuid.UUID, status string) error
}

type GormReviewRepository struct {
	db *gorm.DB
}

func (r *GormReviewRepository) Create(ctx context.Context, review *models.Review) error {
	return r.db.WithContext(ctx).Create(review).Error
}

func (r *GormReviewRepository) ListByCompany(ctx context.Context, companyID uuid.UUID) ([]models.Review, error) {
	var reviews []models.Review
	if err := r.db.WithContext(ctx).
		Preload("User").
		Where("company_id = ?", companyID).
		Order("created_at DESC").
		Find(&reviews).Error; err != nil {
		return nil, err
	}
	return reviews, nil
}

func (r *GormReviewRepository) ListSuspicious(ctx context.Context) ([]models.Review, error) {
	var reviews []models.Review
	if err := r.db.WithContext(ctx).Where("suspicious = ?", true).Preload("ValidationResult").Find(&reviews).Error; err != nil {
		return nil, err
	}
	return reviews, nil
}

func (r *GormReviewRepository) Respond(ctx context.Context, id uuid.UUID, status string) error {
	return r.db.WithContext(ctx).Model(&models.Review{}).Where("id = ?", id).Update("status", status).Error
}
