package models

// OrderItem stores each product in an order with the exact price paid at purchase time.
type OrderItem struct {
	ID         uint  `gorm:"primaryKey"`
	OrderID    uint  `gorm:"index"`
	ProductID  uint  `gorm:"index"`
	Quantity   int   // must be > 0
	PriceCents int64 // price per unit at time of purchase, in cents
}
