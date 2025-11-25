package services

import (
	"crowdreview/config"
	"crowdreview/internal/repository"
	"crowdreview/internal/validation"
	"github.com/redis/go-redis/v9"
)

// Services aggregates service layer dependencies.
type Services struct {
	Auth   AuthService
	Company CompanyService
	Review ReviewService
	Admin  AdminService
}

// NewServices wires concrete service implementations.
func NewServices(cfg config.Config, repos repository.Repositories, rdb *redis.Client, worker *validation.FraudWorker) Services {
	auth := &DefaultAuthService{Users: repos.User, Config: cfg}
	company := &DefaultCompanyService{Companies: repos.Company}
	review := &DefaultReviewService{
		Reviews:     repos.Review,
		Companies:   repos.Company,
		Worker:      worker,
		RateLimiter: rdb,
		Config:      cfg,
	}
	admin := &DefaultAdminService{Reviews: repos.Review, Validation: repos.Validation, DB: repos.DB}

	return Services{
		Auth:    auth,
		Company: company,
		Review:  review,
		Admin:   admin,
	}
}
