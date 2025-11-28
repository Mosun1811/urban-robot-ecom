// This file is assigned to Jasmine, Qas, Yasin ...
// review_handlers.go
//
// PURPOSE:
// - HTTP handlers for product reviews and ratings.
// - Handles creating/updating reviews and listing reviews for a product.
//
// EPICS & USER STORIES:
// - Epic 6: Reviews & Ratings
//   - User Story 6.1: Submit Product Review
//     (POST /api/v1/products/{product_id}/reviews)
//   - User Story 6.2: View Product Reviews
//     (GET /api/v1/products/{product_id}/reviews)
//
// ENDPOINTS (to be implemented here):
// - POST /api/v1/products/{product_id}/reviews
//   - Requires authenticated user.
//   - Accepts rating (1–5) + optional text.
//   - Enforces basic rate limiting (e.g. max 5 reviews/user/min).
//   - If user has already reviewed this product, update existing review.
//
// - GET /api/v1/products/{product_id}/reviews
//   - Public endpoint.
//   - Returns list of reviews sorted by newest first.
//   - Each review should include text, rating, and reviewer display name.
//
// NOTES FOR TEAM:
// - After create/update/delete, service must recalculate product average_rating & review_count.

// What I have done below is just to build so that everything compiles and you'll be able to clone have working code
// Only thing you'd need to do is to write the logic

package handlers

import (
	"futuremarket/service"
	"net/http"

	
)

// ReviewHandler manages reviews and ratings (Epic 6).
type ReviewHandler struct {
	Service service.ReviewService
}

// GET /api/v1/products/{id}/reviews
func (h *ReviewHandler) ListReviews(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("TODO: implement list reviews for a product"))
}

// POST /api/v1/products/{id}/reviews
func (h *ReviewHandler) CreateOrUpdateReview(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("TODO: implement create or update review, recalc average rating"))
}
