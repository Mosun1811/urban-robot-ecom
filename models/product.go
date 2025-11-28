package models

import (
	"time"

	"gorm.io/gorm"
)

// Product is an item that can be listed, searched and bought.
type Product struct {
	gorm.Model
	Name        string `gorm:"size:255"`
	Description string `gorm:"type:text"`
	Category    string `gorm:"size:100"`
	PriceCents  int64  // store price in cents to avoid float issues
	ImageURL    string `gorm:"size:500"`

	// Denormalised rating info (Epic 6.3)
	AverageRating float32 // 1 decimal place logically
	ReviewCount   int64
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
