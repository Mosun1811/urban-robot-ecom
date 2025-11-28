package service

import (
	"futuremarket/models"
	"futuremarket/repository"
)

type UserService struct {
	Repo repository.UserRepo
}

func (us UserService) CreateUser (user *models.User) error {
	return us.Repo.Create(user)
}

func (us UserService) GetUserByEmail (email string) (models.User, error) {
	return us.Repo.GetUserByEmail(email)
}