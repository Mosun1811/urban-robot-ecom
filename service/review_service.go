package service

import (
	"errors"
	"math"

	"futuremarket/models"
	"futuremarket/repository"

	"gorm.io/gorm"
)

// ReviewService contains business logic for Epic 6 (reviews & ratings).
// It uses a ReviewRepo for DB access.
type ReviewService struct {
	Repo repository.ReviewRepo
}

// ListReviews returns review rows + user display names for a product.
func (s *ReviewService) ListReviews(productID uint) ([]repository.ReviewWithUser, error) {
	return s.Repo.ListReviewsWithUser(productID)
}

// CreateOrUpdateReview either updates an existing review for this user+product,
// or creates a new one. It then recalculates average_rating and review_count.
//
// Returns:
//   - the saved review
//   - created = true if it was newly created
//   - error if anything failed
func (s *ReviewService) CreateOrUpdateReview(userID, productID uint, rating int, text string) (*models.Review, bool, error) {
	// Try to find existing review.
	review, err := s.Repo.FindByUserAndProduct(userID, productID)

	if err == nil {
		// Found existing review, update.
		review.Rating = rating
		review.Text = text

		if saveErr := s.Repo.UpdateReview(review); saveErr != nil {
			return nil, false, saveErr
		}

		if recalcErr := s.recalcProductRating(productID); recalcErr != nil {
			return nil, false, recalcErr
		}

		return review, false, nil
	}

	// If not found, create new review.
	if errors.Is(err, gorm.ErrRecordNotFound) {
		newReview := &models.Review{
			ProductID: productID,
			UserID:    userID,
			Rating:    rating,
			Text:      text,
		}

		if createErr := s.Repo.CreateReview(newReview); createErr != nil {
			return nil, false, createErr
		}

		if recalcErr := s.recalcProductRating(productID); recalcErr != nil {
			return nil, false, recalcErr
		}

		return newReview, true, nil
	}

	// Any other DB error.
	return nil, false, err
}

// recalcProductRating recalculates average_rating (1 decimal place)
// and review_count and stores them on the products table.
func (s *ReviewService) recalcProductRating(productID uint) error {
	avg, count, err := s.Repo.CalculateRatingStats(productID)
	if err != nil {
		return err
	}

	avgRounded := math.Round(avg*10) / 10

	return s.Repo.UpdateProductRating(productID, avgRounded, count)

	
}

type PaginatedReviews struct {
    Reviews []models.Review `json:"reviews"`
    Meta    PaginationMeta   `json:"meta"`
}

func (s ReviewService) ListReviewsPaginated(productID uint, page, limit int) (PaginatedReviews, error) {

    reviews, total, err := s.Repo.ListReviewsPaginated(productID, page, limit)
    if err != nil {
        return PaginatedReviews{}, err
    }

    totalPages := int(math.Ceil(float64(total) / float64(limit)))

    meta := PaginationMeta{
        TotalItems:  total,
        TotalPages:  totalPages,
        CurrentPage: page,
        Limit:       limit,
    }

    return PaginatedReviews{
        Reviews: reviews,
        Meta:    meta,
    }, nil
}
