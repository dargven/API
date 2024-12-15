package userService

import (
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
	return user == nil, nil // true, если пользователь не найден
}
