package models

import "gorm.io/gorm"

// Stock tracks how many units are available for a product.
type Stock struct {
	gorm.Model
	ProductID uint `gorm:"uniqueIndex"` // stock is 1-to-1 with product
	Quantity  int
}
