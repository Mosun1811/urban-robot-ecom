package repository

import (
    "futuremarket/models"
    "gorm.io/gorm"
)

type CartRepo struct {
    DB *gorm.DB
}

func (r CartRepo) GetOrCreateCart(userID uint) (*models.Cart, error) {
    var cart models.Cart

    err := r.DB.Where("user_id = ?", userID).First(&cart).Error

    if err == gorm.ErrRecordNotFound {
        cart = models.Cart{UserID: userID}
        if err := r.DB.Create(&cart).Error; err != nil {
            return nil, err
        }
    } else if err != nil {
        return nil, err
    }

    return &cart, nil
}

func (r CartRepo) FindCartItems(cartID uint) ([]models.CartItem, error) {
    var items []models.CartItem
    err := r.DB.Preload("Product").Where("cart_id = ?", cartID).Find(&items).Error
    return items, err
}

func (r CartRepo) AddOrIncreaseItem(cartID, productID uint) error {
    var item models.CartItem

    err := r.DB.Where("cart_id = ? AND product_id = ?", cartID, productID).
        First(&item).Error

    if err == gorm.ErrRecordNotFound {
        item = models.CartItem{
            CartID:    cartID,
            ProductID: productID,
            Quantity:  1,
        }
        return r.DB.Create(&item).Error
    }

    return r.DB.Model(&item).Update("quantity", item.Quantity+1).Error
}

func (r CartRepo) UpdateItemQuantity(cartID, productID uint, qty int) error {
    return r.DB.Model(&models.CartItem{}).
        Where("cart_id = ? AND product_id = ?", cartID, productID).
        Update("quantity", qty).Error
}

func (r CartRepo) RemoveItem(cartID, productID uint) error {
    return r.DB.Where("cart_id = ? AND product_id = ?", cartID, productID).
        Delete(&models.CartItem{}).Error
}
