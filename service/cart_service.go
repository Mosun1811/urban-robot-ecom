package service

import (
	"errors"
	"futuremarket/repository"
)

type CartService struct {
	Repo        repository.CartRepo
	ProductRepo repository.ProductRepo
}

// ADD TO CART
func (s CartService) AddToCart(userID uint, productID uint) error {
	cart, err := s.Repo.GetOrCreateCart(userID)
	if err != nil {
		return err
	}

	// Ensure product exists
	_, err = s.ProductRepo.GetProductByID(productID)
	if err != nil {
		return errors.New("product not found")
	}

	// Check stock table
	stock, err := s.ProductRepo.GetStockByProductID(productID)
	if err != nil {
		return errors.New("stock record missing")
	}
	if stock.Quantity <= 0 {
		return errors.New("product out of stock")
	}

	return s.Repo.AddOrIncreaseItem(cart.ID, productID)
}

// VIEW CART
func (s CartService) GetCart(userID uint) (map[string]any, error) {
	cart, err := s.Repo.GetOrCreateCart(userID)
	if err != nil {
		return nil, err
	}

	items, err := s.Repo.FindCartItems(cart.ID)
	if err != nil {
		return nil, err
	}

	var total int64
	for _, item := range items {
		total += int64(item.Quantity) * item.Product.PriceCents
	}

	return map[string]any{
		"cart_id": cart.ID,
		"items":   items,
		"total":   total,
	}, nil
}

// UPDATE QUANTITY
func (s CartService) UpdateQuantity(userID uint, productID uint, qty int) error {
	if qty <= 0 {
		return errors.New("quantity must be > 0")
	}

	cart, err := s.Repo.GetOrCreateCart(userID)
	if err != nil {
		return err
	}

	_, err = s.ProductRepo.GetProductByID(productID)
	if err != nil {
		return errors.New("product not found")
	}

	stock, err := s.ProductRepo.GetStockByProductID(productID)
	if err != nil {
		return errors.New("stock record missing")
	}

	if qty > stock.Quantity {
		return errors.New("quantity exceeds stock")
	}

	return s.Repo.UpdateItemQuantity(cart.ID, productID, qty)
}

// REMOVE ITEM
func (s CartService) RemoveItem(userID uint, productID uint) error {
	cart, err := s.Repo.GetOrCreateCart(userID)
	if err != nil {
		return err
	}

	return s.Repo.RemoveItem(cart.ID, productID)
}
