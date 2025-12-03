package models

import "gorm.io/gorm"

// TokenBlacklist stores invalidated JWT tokens
type TokenBlacklist struct {
	gorm.Model
	Token string `gorm:"uniqueIndex"`
}
