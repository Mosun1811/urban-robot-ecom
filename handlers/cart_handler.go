package handlers

import (
	"encoding/json"
	"futuremarket/middleware"
	"futuremarket/service"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type CartHandler struct {
	Service service.CartService
}

func getUserIDFromContext(r *http.Request) uint {
	val := r.Context().Value(middleware.ContextUserID).(int)
	return uint(val)
}

// ------------------------------------------------------------
// GET CART
// ------------------------------------------------------------
func (h *CartHandler) GetCart(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)

	resp, err := h.Service.GetCart(userID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(resp)
}

// ------------------------------------------------------------
// ADD TO CART
// ------------------------------------------------------------
func (h *CartHandler) AddToCart(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)

	var body struct {
		ProductID uint `json:"product_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid json", 400)
		return
	}

	err := h.Service.AddToCart(userID, body.ProductID)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	w.WriteHeader(201)
	json.NewEncoder(w).Encode(map[string]string{"message": "added to cart"})
}

// ------------------------------------------------------------
// UPDATE QUANTITY
// ------------------------------------------------------------
func (h *CartHandler) UpdateCartItem(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)

	pid, _ := strconv.Atoi(mux.Vars(r)["product_id"])

	var body struct {
		Quantity int `json:"quantity"`
	}

	json.NewDecoder(r.Body).Decode(&body)

	err := h.Service.UpdateQuantity(userID, uint(pid), body.Quantity)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "quantity updated"})
}

// ------------------------------------------------------------
// REMOVE ITEM
// ------------------------------------------------------------
func (h *CartHandler) RemoveCartItem(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	pid, _ := strconv.Atoi(mux.Vars(r)["product_id"])

	err := h.Service.RemoveItem(userID, uint(pid))
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	w.WriteHeader(204)
}
