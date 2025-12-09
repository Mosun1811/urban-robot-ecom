package service

import (
	"errors"
	"math"

	"futuremarket/models"
	"futuremarket/repository"
)

type ProductService struct {
	Repo repository.ProductRepo
}

type ProductListResponse struct {
	Products []models.Product `json:"products"`
	Meta     PaginationMeta   `json:"meta"`
}

type PaginationMeta struct {
	TotalItems  int64 `json:"total_items"`
	TotalPages  int   `json:"total_pages"`
	CurrentPage int   `json:"current_page"`
	Limit       int   `json:"limit"`
}

// CREATE PRODUCT
func (s ProductService) CreateProduct(p *models.Product) error {
	if p.Name == "" || p.PriceCents <= 0 {
		return errors.New("invalid product fields")
	}
	return s.Repo.CreateProduct(p)
}

// UPDATE PRODUCT
func (s ProductService) UpdateProduct(id uint, updateData *models.Product) (models.Product, error) {
	existing, err := s.Repo.GetProductByID(id)
	if err != nil {
		return models.Product{}, errors.New("product not found")
	}

	if updateData.Name != "" {
		existing.Name = updateData.Name
	}
	if updateData.Description != "" {
		existing.Description = updateData.Description
	}
	if updateData.Category != "" {
		existing.Category = updateData.Category
	}
	if updateData.PriceCents > 0 {
		existing.PriceCents = updateData.PriceCents
	}
	if updateData.ImageURL != "" {
		existing.ImageURL = updateData.ImageURL
	}

	// Only update stock if intentionally set
	if updateData.Stock != 0 {
		existing.Stock = updateData.Stock
	}

	err = s.Repo.UpdateProduct(&existing)
	return existing, err
}

// LIST WITH FILTERS
func (s ProductService) ListProductsWithFilters(
	page int,
	limit int,
	minPrice *int64,
	maxPrice *int64,
	category *string,
) (ProductListResponse, error) {

	if limit <= 0 {
		limit = 10
	}

	products, totalItems, err := s.Repo.ListProductsFiltered(page, limit, minPrice, maxPrice, category)
	if err != nil {
		return ProductListResponse{}, err
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(limit)))

	meta := PaginationMeta{
		TotalItems:  totalItems,
		TotalPages:  totalPages,
		CurrentPage: page,
		Limit:       limit,
	}

	return ProductListResponse{
		Products: products,
		Meta:     meta,
	}, nil
}

func (s ProductService) GetProductByID(id uint) (models.Product, error) {
	return s.Repo.GetProductByID(id)
}
