package models

import "time"

// User is basically either a customer or admin in the system.
//
// Roles:
//   - "customer"
//   - "admin"
type User struct {
	ID           uint   `gorm:"primaryKey"`
	Name         string `gorm:"size:100"`
	Email        string `gorm:"size:255;uniqueIndex"`
	PasswordHash string `gorm:"size:255"`
	Role         string `gorm:"size:20"` // "customer" or "admin"
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
