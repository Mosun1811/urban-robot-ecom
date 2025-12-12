package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"futuremarket/middleware"
	"futuremarket/service"
)

type ReviewHandler struct {
	Service service.ReviewService
}

type reviewRequest struct {
	Rating int    `json:"rating"`
	Text   string `json:"text"`
}

// Utility: parse int query param with default
func parseQueryInt(r *http.Request, key string, defaultVal int) int {
	valStr := r.URL.Query().Get(key)
	if valStr == "" {
		return defaultVal
	}
	n, err := strconv.Atoi(valStr)
	if err != nil || n <= 0 {
		return defaultVal
	}
	return n
}

// -----------------------------------------------------------
// GET /api/v1/products/{id}/reviews   (PUBLIC + PAGINATION)
// -----------------------------------------------------------
func (h *ReviewHandler) ListReviews(w http.ResponseWriter, r *http.Request) {

	page := parseQueryInt(r, "page", 1)
	limit := parseQueryInt(r, "limit", 10)

	idStr := mux.Vars(r)["id"]
	productID, err := strconv.Atoi(idStr)
	if err != nil || productID < 1 {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}

	result, err := h.Service.ListReviewsPaginated(uint(productID), page, limit)
	if err != nil {
		http.Error(w, "unable to load reviews", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(result)
}

// -----------------------------------------------------------
// POST /api/v1/products/{id}/reviews   (AUTH REQUIRED)
// -----------------------------------------------------------
func (h *ReviewHandler) CreateOrUpdateReview(w http.ResponseWriter, r *http.Request) {

	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := mux.Vars(r)["id"]
	productID, err := strconv.Atoi(idStr)
	if err != nil || productID < 1 {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}

	var req reviewRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	if req.Rating < 1 || req.Rating > 5 {
		http.Error(w, "rating must be 1–5", http.StatusBadRequest)
		return
	}

	review, created, err := h.Service.CreateOrUpdateReview(
		userID,
		uint(productID),
		req.Rating,
		req.Text,
	)

	if err != nil {
		http.Error(w, "failed to save review", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if created {
		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	json.NewEncoder(w).Encode(review)
}
