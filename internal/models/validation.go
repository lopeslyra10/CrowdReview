package models

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// ReviewValidationResult stores fraud engine output.
type ReviewValidationResult struct {
	Base
	ReviewID uuid.UUID `gorm:"type:uuid;uniqueIndex"`
	Review   Review    `gorm:"constraint:OnDelete:CASCADE"`
	Score    float64   `gorm:"index"` // 0-100 confidence
	Outcome  string    `gorm:"type:varchar(30)"`
	Checks   datatypes.JSONMap `gorm:"type:jsonb;default:'{}'::jsonb"`
	Signals []FraudSignal `gorm:"foreignKey:ValidationResultID;constraint:OnDelete:CASCADE"`
}

// FraudSignal captures individual rule hits.
type FraudSignal struct {
	Base
	ValidationResultID uuid.UUID `gorm:"type:uuid;index"`
	Type               string    `gorm:"index"`
	Severity           string    `gorm:"type:varchar(10);index"` // low/med/high
	Details            datatypes.JSONMap `gorm:"type:jsonb;default:'{}'::jsonb"`
}
