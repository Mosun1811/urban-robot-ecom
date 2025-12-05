package service

import (
	"futuremarket/repository"
	
)

type CartService struct {
	Repo repository.CartRepo
}