package repository

import (
	"context"

	"crowdreview/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CompanyRepository stores company data.
type CompanyRepository interface {
	Create(ctx context.Context, company *models.Company) error
	Update(ctx context.Context, company *models.Company) error
	List(ctx context.Context) ([]models.Company, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Company, error)
}

type GormCompanyRepository struct {
	db *gorm.DB
}

func (r *GormCompanyRepository) Create(ctx context.Context, company *models.Company) error {
	return r.db.WithContext(ctx).Create(company).Error
}

func (r *GormCompanyRepository) Update(ctx context.Context, company *models.Company) error {
	return r.db.WithContext(ctx).Save(company).Error
}

func (r *GormCompanyRepository) List(ctx context.Context) ([]models.Company, error) {
	var companies []models.Company
	if err := r.db.WithContext(ctx).Preload("Reviews").Find(&companies).Error; err != nil {
		return nil, err
	}
	return companies, nil
}

func (r *GormCompanyRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Company, error) {
	var company models.Company
	if err := r.db.WithContext(ctx).Preload("Reviews").First(&company, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &company, nil
}
