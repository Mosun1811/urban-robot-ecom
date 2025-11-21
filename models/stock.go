package models

// Stock tracks how many units are available for a given product.
type Stock struct {
	ID        uint `gorm:"primaryKey"`
	ProductID uint `gorm:"index"`
	Quantity  int  // must be >= 0
}
