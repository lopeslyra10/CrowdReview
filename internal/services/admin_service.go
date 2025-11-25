package services

import (
	"context"

	"crowdreview/internal/models"
	"crowdreview/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Insights aggregates dashboard metrics.
type Insights struct {
	TotalReviews     int64
	SuspiciousCount  int64
	AverageScore     float64
	CompaniesTracked int64
}

// AdminService exposes admin-only operations.
type AdminService interface {
	GetInsights(ctx context.Context) (Insights, error)
	ListSuspicious(ctx context.Context) ([]models.Review, error)
	Respond(ctx context.Context, reviewID string, status string) error
}

type DefaultAdminService struct {
	Reviews    repository.ReviewRepository
	Validation repository.ValidationRepository
	DB         *gorm.DB
}

func (s *DefaultAdminService) GetInsights(ctx context.Context) (Insights, error) {
	// Simple aggregate queries; defer to DB count to keep example concise.
	var total, suspicious, companies int64
	var avg float64
	s.DB.WithContext(ctx).Model(&models.Review{}).Count(&total)
	s.DB.WithContext(ctx).Model(&models.Review{}).Where("suspicious = ?", true).Count(&suspicious)
	s.DB.WithContext(ctx).Model(&models.Company{}).Count(&companies)
	s.DB.WithContext(ctx).Model(&models.Review{}).Select("COALESCE(avg(rating),0)").Scan(&avg)

	return Insights{
		TotalReviews:     total,
		SuspiciousCount:  suspicious,
		AverageScore:     avg,
		CompaniesTracked: companies,
	}, nil
}

func (s *DefaultAdminService) ListSuspicious(ctx context.Context) ([]models.Review, error) {
	return s.Reviews.ListSuspicious(ctx)
}

func (s *DefaultAdminService) Respond(ctx context.Context, reviewID string, status string) error {
	id, err := uuid.Parse(reviewID)
	if err != nil {
		return err
	}
	return s.Reviews.Respond(ctx, id, status)
}
