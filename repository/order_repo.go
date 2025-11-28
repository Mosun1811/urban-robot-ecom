package repository

import (
	

	"gorm.io/gorm"
)


type OrderRepo struct {
	DB *gorm.DB
}