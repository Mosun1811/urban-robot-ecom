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

// ListOrdersPaginated returns orders with pagination support.
func (r OrderRepo) ListOrdersPaginated(
    userID uint,
    page int,
    limit int,
) ([]models.Order, int64, error) {

    var orders []models.Order
    var total int64

    query := r.DB.Model(&models.Order{}).Where("user_id = ?", userID)

    // Count total results
    if err := query.Count(&total).Error; err != nil {
        return nil, 0, err
    }

    offset := (page - 1) * limit

    // Fetch paginated rows
    err := query.
        Preload("Items").
        Order("created_at DESC").
        Offset(offset).
        Limit(limit).
        Find(&orders).Error

    if err != nil {
        return nil, 0, err
    }

    return orders, total, nil
}

