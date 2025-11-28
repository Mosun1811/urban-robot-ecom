package repository

import (
	

	"gorm.io/gorm"
)


type ProductRepo struct {
	DB *gorm.DB
}