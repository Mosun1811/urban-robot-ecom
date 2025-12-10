package repository

import (
    "futuremarket/models"
    "gorm.io/gorm"
)

type OrderRepo struct {
    DB *gorm.DB
}

func (r OrderRepo) ListOrders(userID uint) ([]models.Order, error) {
    var orders []models.Order
    err := r.DB.Preload("Items").
        Where("user_id = ?", userID).
        Order("created_at DESC").
        Find(&orders).Error

    return orders, err
}
