package models

import (
	"time"

	"gorm.io/gorm"
)

// Review is essentially a customer's review and rating for a product (refer to Epic 6).
type Review struct {
	gorm.Model
	ProductID uint   `gorm:"index"`
	UserID    uint   `gorm:"index"`
	Rating    int    // 1–5
	Text      string `gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
