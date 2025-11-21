package models

// CartItem links a product to a cart with a specific quantity.
type CartItem struct {
	ID        uint `gorm:"primaryKey"`
	CartID    uint `gorm:"index"`
	ProductID uint `gorm:"index"`
	Quantity  int  // must be > 0
}
