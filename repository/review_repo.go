package repository

import (
	

	"gorm.io/gorm"
)


type ReviewRepo struct {
	DB *gorm.DB
}