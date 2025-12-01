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

// ------------------------------------------------------------------
// LIST PRODUCTS WITH PAGINATION + FILTERING
// ------------------------------------------------------------------
func (r ProductRepo) ListProductsFiltered(page, limit int, minPrice, maxPrice *int64, category *string) ([]models.Product, int64, error) {
	var products []models.Product
	var totalItems int64

	query := r.DB.Model(&models.Product{})

	// Apply filters dynamically
	if minPrice != nil {
		query = query.Where("price_cents >= ?", *minPrice)
	}

	if maxPrice != nil {
		query = query.Where("price_cents <= ?", *maxPrice)
	}

	if category != nil && *category != "" {
		query = query.Where("LOWER(category) = LOWER(?)", *category)
	}

	// Count rows AFTER filters
	if err := query.Count(&totalItems).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit

	// Apply pagination
	err := query.
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&products).Error

	if err != nil {
		return nil, 0, err
	}

	return products, totalItems, nil
}

// ------------------------------------------------------------------
// GET SINGLE PRODUCT BY ID
// ------------------------------------------------------------------
func (r ProductRepo) GetProductByID(id uint) (models.Product, error) {
	var product models.Product
	err := r.DB.First(&product, id).Error

	return product, err
}
