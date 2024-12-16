package userService

import (
	"API/internal/models/user"
	repo "API/repositories/userRepository"
	"errors"
	"gorm.io/gorm"
)

type UserService struct {
	repo *repo.UserRepository
}

func NewUserService(repo *repo.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) IsEmailUnique(email string) (bool, error) {
	user, err := s.repo.IsEmailUnique(email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, err // ошибка при запросе
	}
	return user == false, nil // true, если пользователь не найден
}

func (s *UserService) CreateUser(req user.CreateUserRequest) (*user.User, error) {
	newUser := &user.User{
		Email:    req.Email,
		Name:     req.Name,
		Password: req.Password, // Здесь можно добавить хэширование
	}

	if err := s.Repo.Create(newUser); err != nil {
		return nil, err
	}
	return newUser, nil
}
