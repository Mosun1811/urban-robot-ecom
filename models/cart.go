package models

import (
"time"

"gorm.io/gorm"
)

// this one user's basket...
type Cart struct {
gorm.Model
UserID    uint `gorm:"index"`
CreatedAt time.Time
UpdatedAt time.Time
}