package services

import (
	"context"
	"errors"
	"time"

	"crowdreview/config"
	"crowdreview/internal/models"
	"crowdreview/internal/repository"
	"crowdreview/internal/validation"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// ReviewService orchestrates review creation and retrieval.
type ReviewService interface {
	Create(ctx context.Context, userID uuid.UUID, companyID uuid.UUID, input CreateReviewInput) (*models.Review, error)
	ListByCompany(ctx context.Context, companyID uuid.UUID) ([]models.Review, error)
}

// CreateReviewInput is DTO for new reviews.
type CreateReviewInput struct {
	Rating      int
	Title       string
	Content     string
	IPAddress   string
	GeoLocation string
}

type DefaultReviewService struct {
	Reviews     repository.ReviewRepository
	Companies   repository.CompanyRepository
	Worker      *validation.FraudWorker
	RateLimiter *redis.Client
	Config      config.Config
}

func (s *DefaultReviewService) Create(ctx context.Context, userID uuid.UUID, companyID uuid.UUID, input CreateReviewInput) (*models.Review, error) {
	if input.Rating < 1 || input.Rating > 5 {
		return nil, errors.New("rating must be between 1 and 5")
	}

	// Ensure company exists
	if _, err := s.Companies.GetByID(ctx, companyID); err != nil {
		return nil, err
	}

	review := &models.Review{
		UserID:      userID,
		CompanyID:   companyID,
		Rating:      input.Rating,
		Title:       input.Title,
		Content:     input.Content,
		IPAddress:   input.IPAddress,
		GeoLocation: input.GeoLocation,
		Status:      "pending",
	}

	if err := s.Reviews.Create(ctx, review); err != nil {
		return nil, err
	}

	// Enqueue background validation
	s.Worker.Enqueue(*review)

	return review, nil
}

func (s *DefaultReviewService) ListByCompany(ctx context.Context, companyID uuid.UUID) ([]models.Review, error) {
	return s.Reviews.ListByCompany(ctx, companyID)
}
