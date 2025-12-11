package service

import (
	"errors"
	"futuremarket/models"
	"futuremarket/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	Repo repository.UserRepo
}

func (us UserService) CreateUser(user *models.User) error {
	return us.Repo.Create(user)
}

func (us UserService) GetUserByEmail(email string) (models.User, error) {
	return us.Repo.GetUserByEmail(email)
}

//
// ===============================
// REGISTER USER (used by AuthHandler)
// ===============================
//

func (s UserService) RegisterUser(name, email, password, role string) (models.User, error) {

	// 1) VALIDATION
	if err := ValidateName(name); err != nil {
		return models.User{}, err
	}
	if err := ValidateEmail(email); err != nil {
		return models.User{}, err
	}
	if err := ValidatePassword(password); err != nil {
		return models.User{}, err
	}

	// 2) Check duplicate email
	_, err := s.Repo.GetUserByEmail(email)
	if err == nil {
		return models.User{}, errors.New("email already registered")
	}

	// 3) Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, err
	}

	// 4) Create user
	newUser := models.User{
		Name:         name,
		Email:        email,
		PasswordHash: string(hashed),
		Role:         role,
	}

	if err := s.Repo.Create(&newUser); err != nil {
		return models.User{}, err
	}

	return newUser, nil
}
