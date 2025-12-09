package models

import "gorm.io/gorm"

// CartItem links a product to a cart with a quantity.
type CartItem struct {
	gorm.Model
	CartID    uint    `gorm:"index"`
	ProductID uint    `gorm:"index"`
	Quantity  int     // must be > 0
	Product   Product `gorm:"foreignKey:ProductID"` // preload support
}
