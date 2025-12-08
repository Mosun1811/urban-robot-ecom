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

// Create a new product
func (r ProductRepo) CreateProduct(product *models.Product) error {
	return r.DB.Create(product).Error
}

// Update existing product by ID
func (r ProductRepo) UpdateProduct(product *models.Product) error {
	return r.DB.Save(product).Error
}

// Fetch a single product
func (r ProductRepo) GetProductByID(id uint) (models.Product, error) {
	var product models.Product
	err := r.DB.First(&product, id).Error
	return product, err
}

func (r ProductRepo) ListProductsFiltered(
	page int,
	limit int,
	minPrice *int64,
	maxPrice *int64,
	category *string,
) ([]models.Product, int64, error) {

	var products []models.Product
	var totalItems int64

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

	// Count after filters
	if err := query.Count(&totalItems).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit

	// Get paginated results
	err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&products).Error

	if err != nil {
		return nil, 0, err
	}

	return products, totalItems, nil
}
