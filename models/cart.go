package models

import "time"

// this one user's basket...
type Cart struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint `gorm:"index"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
