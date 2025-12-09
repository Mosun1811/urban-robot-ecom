package service

import (
	"errors"
	"fmt"

	"futuremarket/models"
	"futuremarket/repository"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OrderService struct {
	OrderRepo   repository.OrderRepo
	CartRepo    repository.CartRepo
	ProductRepo repository.ProductRepo
}

func (s OrderService) Checkout(userID uint) error {

	db := s.OrderRepo.DB
	if db == nil {
		return errors.New("order repository db is nil")
	}

	return db.Transaction(func(tx *gorm.DB) error {

		// ----------------------------------------------------
		// 1) Load user's cart
		// ----------------------------------------------------
		cart, err := s.CartRepo.GetOrCreateCart(userID)
		if err != nil {
			return err
		}

		items, err := s.CartRepo.FindCartItems(cart.ID)
		if err != nil {
			return err
		}
		if len(items) == 0 {
			return fmt.Errorf("cart is empty")
		}

		// ----------------------------------------------------
		// 2) Stock checks + stock deduction inside TX
		// ----------------------------------------------------
		var total int64 = 0
		orderItems := make([]models.OrderItem, 0, len(items))

		for _, ci := range items {

			// Lock product row
			var product models.Product
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
				Where("id = ?", ci.ProductID).First(&product).Error; err != nil {
				return err
			}

			// Lock stock row
			stock, err := s.ProductRepo.GetStockLocked(tx, ci.ProductID)
			if err != nil {
				return fmt.Errorf("missing stock record for product %d", ci.ProductID)
			}

			// Insufficient stock?
			if stock.Quantity < ci.Quantity {
				return fmt.Errorf("insufficient stock for product %d", ci.ProductID)
			}

			// Deduct stock
			newQty := stock.Quantity - ci.Quantity
			if err := tx.Model(&models.Stock{}).
				Where("id = ?", stock.ID).
				Update("quantity", newQty).Error; err != nil {
				return err
			}

			// Build OrderItem record
			orderItems = append(orderItems, models.OrderItem{
				ProductID:  ci.ProductID,
				Quantity:   ci.Quantity,
				PriceCents: product.PriceCents,
			})

			total += int64(ci.Quantity) * product.PriceCents
		}

		// ----------------------------------------------------
		// 3) Create Order
		// ----------------------------------------------------
		order := models.Order{
			UserID: userID,
			Status: "Pending",
			Total:  total,
		}
		if err := tx.Create(&order).Error; err != nil {
			return err
		}

		// Attach orderID
		for i := range orderItems {
			orderItems[i].OrderID = order.ID
		}

		// Insert order items
		if err := tx.Create(&orderItems).Error; err != nil {
			return err
		}

		// ----------------------------------------------------
		// 4) Clear cart
		// ----------------------------------------------------
		if err := tx.Where("cart_id = ?", cart.ID).Delete(&models.CartItem{}).Error; err != nil {
			return err
		}

		return nil
	})
}

// ------------------------------------------------------------
// LIST ORDERS
// ------------------------------------------------------------
func (s OrderService) ListOrders(userID uint) ([]models.Order, error) {

	db := s.OrderRepo.DB
	if db == nil {
		return nil, errors.New("order repo db is nil")
	}

	var orders []models.Order

	err := db.
		Preload("Items").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&orders).Error

	return orders, err
}
