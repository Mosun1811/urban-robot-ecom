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
	
)

// ProductHandler manages product listing, search and admin product management.
type ProductHandler struct {
	Service service.ProductService
}

// GET /api/v1/products
func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	// Parse pagination params: ?page=2&limit=50
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page := 1
	limit := 20 // default limit

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	// Call service
	result, err := h.Service.ListProducts(page, limit)
	if err != nil {
		http.Error(w, "failed to fetch products", http.StatusInternalServerError)
		return
	}

	// Send JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// GET /api/v1/products/{id}
func (h *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("TODO: implement get single product by ID"))
}

// POST /api/v1/admin/products  (via admin routes)
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("TODO: implement admin create product"))
}

// PATCH /api/v1/admin/products/{id}
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("TODO: implement admin update product"))
}
