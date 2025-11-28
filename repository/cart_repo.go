package repository

import (
	

	"gorm.io/gorm"
)


type CartRepo struct {
	DB *gorm.DB
}