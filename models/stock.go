package models

import "gorm.io/gorm"

// Stock tracks how many units are available for a given product.
type Stock struct {
	gorm.Model
	ProductID uint `gorm:"index"`
	Quantity  int  // must be >= 0
}
