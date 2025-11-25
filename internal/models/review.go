package models

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// Review represents a user review awaiting validation.
type Review struct {
	Base
	UserID             uuid.UUID
	User               User `gorm:"constraint:OnDelete:CASCADE"`
	CompanyID          uuid.UUID `gorm:"index"`
	Company            Company   `gorm:"constraint:OnDelete:CASCADE"`
	Rating             int       `gorm:"check:rating BETWEEN 1 AND 5"`
	Title              string
	Content            string            `gorm:"type:text"`
	IPAddress          string            `gorm:"index"`
	GeoLocation        string
	Status             string            `gorm:"type:varchar(20);index;default:'pending'"`
	Suspicious         bool              `gorm:"index"`
	ValidationResultID *uuid.UUID
	ValidationResult   *ReviewValidationResult
	Metadata           datatypes.JSONMap `gorm:"type:jsonb;default:'{}'::jsonb"`
}
