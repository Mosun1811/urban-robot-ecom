package handlers

import (
	"encoding/json"
	"futuremarket/middleware"
	"futuremarket/service"
	"net/http"
	"strconv"
)

type OrderHandler struct {
	Service service.OrderService
}

func getUserID(r *http.Request) uint {
	return uint(r.Context().Value(middleware.ContextUserID).(int))
}

func (h *OrderHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)

	if err := h.Service.Checkout(userID); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "order created successfully",
	})
}

func (h *OrderHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)

	orders, err := h.Service.ListOrders(userID)
	if err != nil {
		http.Error(w, "failed to load orders", 500)
		return
	}

	json.NewEncoder(w).Encode(orders)
}

// GET /api/v1/orders/paginated?page=1&limit=10
func (h *OrderHandler) ListOrdersPaginated(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)

	// Read query params
	page := 1
	limit := 20

	if p := r.URL.Query().Get("page"); p != "" {
		if val, err := strconv.Atoi(p); err == nil && val > 0 {
			page = val
		}
	}
	if l := r.URL.Query().Get("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil && val > 0 {
			limit = val
		}
	}

	// Call service
	orders, total, err := h.Service.ListOrdersPaginated(userID, page, limit)
	if err != nil {
		http.Error(w, "failed to load orders", http.StatusInternalServerError)
		return
	}

	// Build response
	response := map[string]any{
		"orders": orders,
		"meta": map[string]any{
			"total_items": total,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
			"page":        page,
			"limit":       limit,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
