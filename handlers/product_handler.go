// This file is assigned to Valencia
// PURPOSE:
// - HTTP handlers for product catalogue and admin product management.
// - Handles listing, filtering, fetching details, and admin create/update.
//
// EPICS & USER STORIES:
// - Epic 2: Product Catalog & Discovery
//   - User Story 2.1: Product Listing with Pagination (GET /api/v1/products)
//   - User Story 2.2: Product Search & Filtering       (GET /api/v1/products with query params)
//   - User Story 2.3: Retrieve Single Product Details  (GET /api/v1/products/{product_id})
// - Epic 5: Administrator Dashboard
//   - User Story 5.1: Product Management (POST/PATCH /api/v1/products)
//
// ENDPOINTS (to be implemented here):
// - GET /api/v1/products
//   - Supports pagination: ?page=&limit=
//   - Supports filters: ?min_price=&max_price=&category=
//   - Returns: list of products + metadata (total_items, total_pages, current_page).
//
// - GET /api/v1/products/{product_id}
//   - Returns full product details.
//   - 404 if not found.
//
// - POST /api/v1/products          (admin only)
//   - Create a new product.
//
// - PATCH /api/v1/products/{id}    (admin only)
//   - Update existing product fields.
//

// What I have done below is just to build so that everything compiles and you'll be able to clone have working code
// Only thing you'd need to do is to write the logic

package handlers

import (
	"encoding/json"
	"futuremarket/service"
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
	"futuremarket/models"


	
)

// ProductHandler manages product listing, search and admin product management.
type ProductHandler struct {
	Service service.ProductService
}

// GET /api/v1/products
func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	// Pagination
	page := 1
	limit := 20

	if p, err := strconv.Atoi(q.Get("page")); err == nil && p > 0 {
		page = p
	}
	if l, err := strconv.Atoi(q.Get("limit")); err == nil && l > 0 {
		limit = l
	}

	// Filters
	var minPricePtr, maxPricePtr *int64
	if v := q.Get("min_price"); v != "" {
		if val, err := strconv.ParseInt(v, 10, 64); err == nil {
			minPricePtr = &val
		}
	}
	if v := q.Get("max_price"); v != "" {
		if val, err := strconv.ParseInt(v, 10, 64); err == nil {
			maxPricePtr = &val
		}
	}

	category := q.Get("category")
	var categoryPtr *string
	if category != "" {
		categoryPtr = &category
	}

	// Call service
	result, err := h.Service.ListProductsWithFilters(page, limit, minPricePtr, maxPricePtr, categoryPtr)
	if err != nil {
		http.Error(w, "failed to fetch products", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}


// GET /api/v1/products/{id}
func (h *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}

	product, err := h.Service.GetProductByID(uint(id))
	if err != nil {
		http.Error(w, "product not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(product)
}


// POST /api/v1/admin/products  (via admin routes)
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Category    string `json:"category"`
		PriceCents  int64  `json:"price_cents"`
		Stock       int64  `json:"stock"` 
		ImageURL    string `json:"image_url"`
	}

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	// Basic validation
	if req.Name == "" || req.PriceCents <= 0 {
		http.Error(w, "name and price_cents are required", http.StatusBadRequest)
		return
	}

	// Build product model
	product := models.Product{
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
		PriceCents:  req.PriceCents,
		Stock:		 req.Stock,
		ImageURL:    req.ImageURL,
		// average rating fields start at zero
	}

	// Call service layer
	err := h.Service.CreateProduct(&product)
	if err != nil {
		http.Error(w, "failed to create product", http.StatusInternalServerError)
		return
	}

	// Success response
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

// PATCH /api/v1/admin/products/{id}
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}


	var req models.Product

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	// Call service to update only provided fields
	updated, err := h.Service.UpdateProduct(uint(id), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}

