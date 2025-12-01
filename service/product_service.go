package service

import (
	"futuremarket/models"
	"futuremarket/repository"
	"math"
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

func (s ProductService) ListProducts(page, limit int) (ProductListResponse, error) {
	products, totalItems, err := s.Repo.ListProducts(page, limit)
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
