package repository

import (
    "futuremarket/models"
    "gorm.io/gorm"
)

type StockRepo struct {
    DB *gorm.DB
}

func (r StockRepo) GetStockForProduct(productID uint) (models.Stock, error) {
    var s models.Stock
    err := r.DB.Where("product_id = ?", productID).First(&s).Error
    return s, err
}

func (r StockRepo) UpdateStock(stock models.Stock) error {
    return r.DB.Save(&stock).Error
}
