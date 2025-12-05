package repository

import (
	"time"

	"futuremarket/models"

	"gorm.io/gorm"
)

// ReviewWithUser is a read model that includes the review
// plus the reviewer's display name from the users table.
type ReviewWithUser struct {
	ID          uint
	ProductID   uint
	UserID      uint
	Rating      int
	Text        string
	CreatedAt   time.Time
	DisplayName string
}

// ReviewRepo wraps all DB operations for reviews and rating aggregates.
type ReviewRepo struct {
	DB *gorm.DB
}

// ListReviewsWithUser returns reviews for a product, joined with users
// and ordered by newest first.
func (r *ReviewRepo) ListReviewsWithUser(productID uint) ([]ReviewWithUser, error) {
	var rows []ReviewWithUser

	err := r.DB.
		Table("reviews").
		Select("reviews.id, reviews.product_id, reviews.user_id, reviews.rating, reviews.text, reviews.created_at, users.name as display_name").
		Joins("JOIN users ON users.id = reviews.user_id").
		Where("reviews.product_id = ?", productID).
		Order("reviews.created_at DESC").
		Scan(&rows).Error

	if err != nil {
		return nil, err
	}

	return rows, nil
}

// FindByUserAndProduct fetches a single review by user and product.
func (r *ReviewRepo) FindByUserAndProduct(userID, productID uint) (*models.Review, error) {
	var review models.Review

	err := r.DB.
		Where("user_id = ? AND product_id = ?", userID, productID).
		First(&review).Error

	if err != nil {
		return nil, err
	}

	return &review, nil
}

// CreateReview inserts a new review row.
func (r *ReviewRepo) CreateReview(review *models.Review) error {
	return r.DB.Create(review).Error
}

// UpdateReview updates an existing review row.
func (r *ReviewRepo) UpdateReview(review *models.Review) error {
	return r.DB.Save(review).Error
}

// CalculateRatingStats returns the average rating and total review count
// for a given product.
func (r *ReviewRepo) CalculateRatingStats(productID uint) (float64, int64, error) {
	var avg float64
	var count int64

	err := r.DB.
		Model(&models.Review{}).
		Where("product_id = ?", productID).
		Select("COALESCE(AVG(rating), 0)").
		Scan(&avg).Error
	if err != nil {
		return 0, 0, err
	}

	err = r.DB.
		Model(&models.Review{}).
		Where("product_id = ?", productID).
		Count(&count).Error
	if err != nil {
		return 0, 0, err
	}

	return avg, count, nil
}

// UpdateProductRating writes the denormalized average_rating and review_count
// back onto the products table.
func (r *ReviewRepo) UpdateProductRating(productID uint, avg float64, count int64) error {
	return r.DB.
		Model(&models.Product{}).
		Where("id = ?", productID).
		Updates(map[string]interface{}{
			"average_rating": avg,
			"review_count":   count,
		}).Error
}
