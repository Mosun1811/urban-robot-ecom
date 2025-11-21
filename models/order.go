package models

import "time"

// Order is basically when a user checks out their cart (please refer to Epic 4 if needed).
type Order struct {
	ID         uint   `gorm:"primaryKey"`
	UserID     uint   `gorm:"index"`
	Status     string `gorm:"size:20"` // e.g. "PENDING", "SHIPPED", "CANCELLED"
	TotalCents int64  // snapshot of total at time of purchase
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
