package repository

import (
	"futuremarket/models"

	"gorm.io/gorm"
)

type ProductRepo struct {
	DB *gorm.DB
}

func NewProductRepo(db *gorm.DB) ProductRepo {
	return ProductRepo{DB: db}
}

// ListProducts retrieves products with pagination.
func (r ProductRepo) ListProducts(page int, limit int) ([]models.Product, int64, error) {
	var products []models.Product
	var totalItems int64

	// Count total product rows
	if err := r.DB.Model(&models.Product{}).Count(&totalItems).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit

	// Fetch paginated products
	err := r.DB.
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&products).Error

	if err != nil {
		return nil, 0, err
	}

	return products, totalItems, nil
}
