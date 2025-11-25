package services

import (
	"context"

	"crowdreview/internal/models"
	"crowdreview/internal/repository"
	"github.com/google/uuid"
)

// CompanyService handles company CRUD.
type CompanyService interface {
	List(ctx context.Context) ([]models.Company, error)
	Get(ctx context.Context, id uuid.UUID) (*models.Company, error)
	Create(ctx context.Context, input models.Company) (*models.Company, error)
	Update(ctx context.Context, id uuid.UUID, input models.Company) (*models.Company, error)
}

type DefaultCompanyService struct {
	Companies repository.CompanyRepository
}

func (s *DefaultCompanyService) List(ctx context.Context) ([]models.Company, error) {
	return s.Companies.List(ctx)
}

func (s *DefaultCompanyService) Get(ctx context.Context, id uuid.UUID) (*models.Company, error) {
	return s.Companies.GetByID(ctx, id)
}

func (s *DefaultCompanyService) Create(ctx context.Context, input models.Company) (*models.Company, error) {
	if err := s.Companies.Create(ctx, &input); err != nil {
		return nil, err
	}
	return &input, nil
}

func (s *DefaultCompanyService) Update(ctx context.Context, id uuid.UUID, input models.Company) (*models.Company, error) {
	company, err := s.Companies.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	company.Name = input.Name
	company.Description = input.Description
	company.Industry = input.Industry
	company.Domain = input.Domain
	company.Website = input.Website
	company.Location = input.Location
	if err := s.Companies.Update(ctx, company); err != nil {
		return nil, err
	}
	return company, nil
}
