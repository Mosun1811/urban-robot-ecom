package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"futuremarket/service"
	"futuremarket/middleware"
)

// ReviewHandler manages reviews and ratings (Epic 6).
// It depends on a ReviewService, which is passed in from main.go.
type ReviewHandler struct {
	Service service.ReviewService
}

// reviewRequest is the expected JSON body for POST /reviews.
type reviewRequest struct {
	Rating int    `json:"rating"`
	Text   string `json:"text"`
}

// reviewResponse is the JSON shape returned by GET /reviews.
type reviewResponse struct {
	ID          uint      `json:"id"`
	ProductID   uint      `json:"product_id"`
	UserID      uint      `json:"user_id"`
	Rating      int       `json:"rating"`
	Text        string    `json:"text"`
	DisplayName string    `json:"display_name"`
	CreatedAt   time.Time `json:"created_at"`
}

// rateLimiter stores recent review timestamps per user for basic rate limiting.
type rateLimiter struct {
	mu   sync.Mutex
	hits map[uint][]time.Time
}

func newRateLimiter() *rateLimiter {
	return &rateLimiter{
		hits: make(map[uint][]time.Time),
	}
}

// allow returns true if the user has made fewer than 5 review requests
// in the last 60 seconds. Otherwise it blocks the request.
func (rl *rateLimiter) allow(userID uint) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-1 * time.Minute)

	history := rl.hits[userID]
	var kept []time.Time

	for _, t := range history {
		if t.After(windowStart) {
			kept = append(kept, t)
		}
	}

	if len(kept) >= 5 {
		rl.hits[userID] = kept
		return false
	}

	kept = append(kept, now)
	rl.hits[userID] = kept
	return true
}

// global in-memory limiter for all review requests.
var reviewLimiter = newRateLimiter()

// getProductIDFromPath extracts the {product_id} from a path like
// /api/v1/products/12/reviews.
func getProductIDFromPath(path string) (uint, error) {
	parts := strings.Split(strings.Trim(path, "/"), "/")

	for i := 0; i < len(parts); i++ {
		if parts[i] == "products" && i+1 < len(parts) {
			idStr := parts[i+1]

			idInt, err := strconv.Atoi(idStr)
			if err != nil || idInt <= 0 {
				return 0, errors.New("invalid product id")
			}

			return uint(idInt), nil
		}
	}

	return 0, errors.New("product id not found")
}

func getUserIDFromContext(r *http.Request) (uint, bool) {
	val := r.Context().Value(middleware.ContextUserID)
	if val == nil {
		return 0, false
	}

	if id, ok := val.(int); ok && id > 0 {
		return uint(id), true
	}

	return 0, false
}

// ListReviews handles:
//   GET /api/v1/products/{id}/reviews
//
// Public endpoint.
// Returns all reviews for a product, newest first,
// including rating, text, and reviewer display name.
func (h *ReviewHandler) ListReviews(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	productID, err := getProductIDFromPath(r.URL.Path)
	if err != nil {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}

	rows, err := h.Service.ListReviews(productID)
	if err != nil {
		http.Error(w, "failed to fetch reviews", http.StatusInternalServerError)
		return
	}

	var resp []reviewResponse

	for _, rr := range rows {
		resp = append(resp, reviewResponse{
			ID:          rr.ID,
			ProductID:   rr.ProductID,
			UserID:      rr.UserID,
			Rating:      rr.Rating,
			Text:        rr.Text,
			DisplayName: rr.DisplayName,
			CreatedAt:   rr.CreatedAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// CreateOrUpdateReview handles:
//   POST /api/v1/products/{id}/reviews
//
// Requires authenticated user.
// Validates rating, enforces rate limit, and calls the ReviewService to
// update or create the review.
func (h *ReviewHandler) CreateOrUpdateReview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	//TEMP for testing: pretend user with ID 1 is logged in
	userID := uint(1)

	productID, err := getProductIDFromPath(r.URL.Path)
	if err != nil {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}

	// Basic in-memory rate limiting: max 5 reviews per user per minute.
	if !reviewLimiter.allow(userID) {
		http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
		return
	}

	var req reviewRequest

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	if req.Rating < 1 || req.Rating > 5 {
		http.Error(w, "rating must be 1 to 5", http.StatusBadRequest)
		return
	}

	review, created, err := h.Service.CreateOrUpdateReview(userID, productID, req.Rating, req.Text)
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
