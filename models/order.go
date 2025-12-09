package models

import (
    "gorm.io/gorm"
    "time"
)

type Order struct {
    gorm.Model

    UserID uint   `gorm:"index"`
    Status string
    Total  int64

    Items []OrderItem `gorm:"foreignKey:OrderID"`

    CreatedAt time.Time
    UpdatedAt time.Time
}
