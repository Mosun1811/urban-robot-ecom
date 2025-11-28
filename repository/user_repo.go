package repository

import (
	"futuremarket/models"
	

	"gorm.io/gorm"
)

type UserRepo struct {
	DB *gorm.DB
}

func (ur UserRepo) Create(user *models.User) error {
	err := ur.DB.Create(&user).Error
	if err != nil {
		return err
	}
	
	return nil
}

func (ur UserRepo) GetUserByEmail(email string) (models.User, error) {
	var existing *models.User
	if err := ur.DB.Where("email = ?", email).First(&existing).Error; err != nil {

		return *existing, err
	}
	
	return *existing, nil
}
