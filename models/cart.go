package models

import "gorm.io/gorm"

type Cart struct {
    gorm.Model
    UserID uint `gorm:"index"`

    Items []CartItem `gorm:"foreignKey:CartID"`
}
