package models

import (
	

	"gorm.io/gorm"
)

// User is basically either a customer or admin in the system.
//
// Roles:
//   - "customer"
//   - "admin"
type User struct {
	gorm.Model
	Name         string `gorm:"size:100"`
	Email        string `gorm:"size:255;uniqueIndex"`
	PasswordHash string `gorm:"size:255"`
	Role         string `gorm:"size:20"` // "customer" or "admin"
	
}
