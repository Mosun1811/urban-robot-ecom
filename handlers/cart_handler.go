// This file is assigned to Jordan, please have a look, let me know if you need anything...
// cart_handlers.go
//
// PURPOSE:
// - HTTP handlers for the shopping cart feature.
// - All endpoints here needs an authenticated user (JWT).
//
// EPICS & USER STORIES:
// - Epic 3: Shopping Cart & Inventory
//   - User Story 3.1: Add to Cart   (POST /api/v1/cart or similar)
//   - User Story 3.2: View Cart      (GET /api/v1/cart)
//   - User Story 3.3: Update Item Quantity (PATCH /api/v1/cart/{product_id})
//   - User Story 3.4: Remove Item   (DELETE /api/v1/cart/{product_id})
//
// ENDPOINTS (to be implemented here):
// - POST /api/v1/cart or /api/v1/cart/{product_id}
//   - Add a product to the authenticated user's cart.
//   - If already exists, increment quantity.
//   - Must check stock > 0 before adding.
//
// - GET /api/v1/cart
//   - Returns all items in the user's cart.
//   - Includes unit prices and dynamically calculated total.
//
// - PATCH /api/v1/cart/{product_id}
//   - Update quantity of a specific product in the cart.
//   - Quantity must be > 0 and <= available stock.
//
// - DELETE /api/v1/cart/{product_id}
//   - Remove an item from the user's cart.
//

// - Cart totals should be calculated using current product prices (no floats: use cents/int).

package handlers

import (
	"futuremarket/service"
	"net/http"

	
)

// CartHandler manages the shopping cart (Epic 3).
type CartHandler struct {
	Service service.CartService
}

// GET /api/v1/cart
func (h *CartHandler) GetCart(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("TODO: implement view cart with total price"))
}

// POST /api/v1/cart
func (h *CartHandler) AddToCart(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("TODO: implement add item to cart (check stock, increment quantity)"))
}

// PATCH /api/v1/cart/{product_id}
func (h *CartHandler) UpdateCartItem(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("TODO: implement update item quantity in cart"))
}

// DELETE /api/v1/cart/{product_id}
func (h *CartHandler) RemoveCartItem(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("TODO: implement remove item from cart"))
}

// What I have done below is just to build so that everything compiles and you'll be able to clone have working code
// Only thing you'd need to do is to write the logic