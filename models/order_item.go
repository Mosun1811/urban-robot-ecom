package models

import "gorm.io/gorm"

type OrderItem struct {
    gorm.Model

    OrderID    uint `gorm:"index"`
    ProductID  uint `gorm:"index"`
    Quantity   int
    PriceCents int64
}
