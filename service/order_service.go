package service

import (
	"futuremarket/repository"
	
)

type OrderService struct {
	Repo repository.OrderRepo
}