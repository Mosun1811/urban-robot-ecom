package handlers

import (
	"encoding/json"
	"futuremarket/middleware"
	"futuremarket/service"
	"net/http"
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
