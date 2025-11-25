package models

import "gorm.io/datatypes"

// Company represents a company that receives reviews.
type Company struct {
	Base
	Name        string            `gorm:"uniqueIndex:idx_companies_name,where:deleted_at IS NULL;not null"`
	Domain      string            `gorm:"uniqueIndex:idx_companies_domain,where:deleted_at IS NULL"`
	Industry    string            `gorm:"index"`
	Location    string
	Description string            `gorm:"type:text"`
	Website     string
	Metrics     datatypes.JSONMap `gorm:"type:jsonb;default:'{}'::jsonb"` // dashboard metrics cache
	Reviews     []Review
}
