// this file is assigned to Jasmine and Hina,

// PURPOSE:
// - HTTP handlers for checkout and order history.
// - Responsible for turning the cart into an order using DB transactions.
//
// EPICS & USER STORIES:
// - Epic 4: Checkout & Orders
//   - User Story 4.1: Place Order    (POST /api/v1/checkout)
//   - User Story 4.2: Order History  (GET /api/v1/orders)
//
// ENDPOINTS (to be implemented here):
// - POST /api/v1/checkout
//   - Uses the authenticated user's cart.
//   - Calls a service that performs an ACID transaction:
//     1) Re-check stock (with row locking if needed).
//     2) Deduct stock from products/stock table.
//     3) Create new order record.
//     4) Move cart_items → order_items.
//     5) Clear the cart.
//   - If ANY step fails, transaction MUST roll back.
//
// - GET /api/v1/orders
//   - Returns list of past orders for the authenticated user.
//   - Includes status (Pending, Shipped, Cancelled).
//

// Ladies just remember this,  Handler just needs::
//   - Gets user_id from JWT
//   - Calls orderService.Checkout(userID)
//   - Returns success/error response.

// What I have done below is just to build so that everything compiles and you'll be able to clone have working code
// Only thing you'd need to do is to write the logic

// What I have done below is just to build so that everything compiles and you'll be able to clone have working code
// Only thing you'd need to do is to write the logic

package handlers

import (
	"futuremarket/service"
	"net/http"

	
)

// OrderHandler manages checkout and order history (Epic 4).
type OrderHandler struct {
	Service service.OrderService
}

// POST /api/v1/checkout
func (h *OrderHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("TODO: implement checkout transaction (cart -> order)"))
}

// GET /api/v1/orders
func (h *OrderHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("TODO: implement list past orders for current user"))
}
