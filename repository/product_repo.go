package repository

import (
	"futuremarket/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProductRepo struct {
	DB *gorm.DB
}

func NewProductRepo(db *gorm.DB) ProductRepo {
	return ProductRepo{DB: db}
}

// Create product
func (r ProductRepo) CreateProduct(product *models.Product) error {
	return r.DB.Create(product).Error
}

// Update product
func (r ProductRepo) UpdateProduct(product *models.Product) error {
	return r.DB.Save(product).Error
}

// Get product
func (r ProductRepo) GetProductByID(id uint) (models.Product, error) {
	var product models.Product
	err := r.DB.First(&product, id).Error
	return product, err
}

// ⭐ REAL STOCK LOOKUP (used by CartService)
func (r ProductRepo) GetStockByProductID(productID uint) (*models.Stock, error) {
	var stock models.Stock
	err := r.DB.Where("product_id = ?", productID).First(&stock).Error
	if err != nil {
		return nil, err
	}
	return &stock, nil
}

// ⭐ LOCKED STOCK LOOKUP (used by OrderService inside transactions)
func (r ProductRepo) GetStockLocked(tx *gorm.DB, productID uint) (*models.Stock, error) {
	var stock models.Stock
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("product_id = ?", productID).
		First(&stock).Error

	if err != nil {
		return nil, err
	}
	return &stock, nil
}

// ListProductsFiltered applies pagination + filtering for Epic 2
func (r ProductRepo) ListProductsFiltered(
	page int,
	limit int,
	minPrice *int64,
	maxPrice *int64,
	category *string,
) ([]models.Product, int64, error) {

	var products []models.Product
	var totalItems int64

	// Start query
	query := r.DB.Model(&models.Product{})

	// Apply filters
	if minPrice != nil {
		query = query.Where("price_cents >= ?", *minPrice)
	}
	if maxPrice != nil {
		query = query.Where("price_cents <= ?", *maxPrice)
	}
	if category != nil {
		query = query.Where("category = ?", *category)
	}

	// Count total after filters
	if err := query.Count(&totalItems).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit

	// Retrieve paginated products
	err := query.
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&products).Error

	if err != nil {
		return nil, 0, err
	}

	return products, totalItems, nil
}
