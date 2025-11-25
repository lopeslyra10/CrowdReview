package repository

import "gorm.io/gorm"

// Repositories aggregates all repo interfaces for easy injection.
type Repositories struct {
	User        UserRepository
	Company     CompanyRepository
	Review      ReviewRepository
	Validation  ValidationRepository
	Achievement AchievementRepository
	DB          *gorm.DB
}

// NewRepositories wires GORM implementations.
func NewRepositories(db *gorm.DB) Repositories {
	return Repositories{
		User:        &GormUserRepository{db},
		Company:     &GormCompanyRepository{db},
		Review:      &GormReviewRepository{db},
		Validation:  &GormValidationRepository{db},
		Achievement: &GormAchievementRepository{db},
		DB:          db,
	}
}
