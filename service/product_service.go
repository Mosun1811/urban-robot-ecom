package service

import (
	"futuremarket/models"
	"futuremarket/repository"
	"math"
)

// --------------------------------------
// RESPONSE STRUCTS (must be above methods)
// --------------------------------------

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

// --------------------------------------
// SERVICE STRUCT
// --------------------------------------

type ProductService struct {
	Repo repository.ProductRepo
}

// ------------------------------------------------------------------
// LIST PRODUCTS WITH FILTERING
// ------------------------------------------------------------------

func (s ProductService) ListProductsWithFilters(
	page, limit int,
	minPrice, maxPrice *int64,
	category *string,
) (ProductListResponse, error) {

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

// ------------------------------------------------------------------
// GET PRODUCT BY ID
// ------------------------------------------------------------------

func (s ProductService) GetProductByID(id uint) (models.Product, error) {
	return s.Repo.GetProductByID(id)
}
