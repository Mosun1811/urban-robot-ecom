package repository

import (
	"futuremarket/models"

	"gorm.io/gorm"
)

type BlacklistRepository struct {
	DB *gorm.DB
}

func NewBlacklistRepository(db *gorm.DB) BlacklistRepository {
	return BlacklistRepository{DB: db}
}

// Add stores a token in the blacklist
func (r BlacklistRepository) Add(token string) error {
	entry := models.TokenBlacklist{Token: token}
	return r.DB.Create(&entry).Error
}

// Exists checks if a token has been blacklisted already
func (r BlacklistRepository) Exists(token string) (bool, error) {
	var entry models.TokenBlacklist
	err := r.DB.Where("token = ?", token).First(&entry).Error

	if err == gorm.ErrRecordNotFound {
		return false, nil
	}

	return err == nil, err
}
