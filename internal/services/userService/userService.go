package userService

import (
	user "API/internal/models/user"
	repo "API/repositories/userRepository"
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	repo *repo.UserRepository
}

func NewUserService(repo *repo.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) IsEmailUnique(ctx context.Context, email string) (bool, error) {
	user, err := s.repo.IsEmailUnique(ctx, email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, err // ошибка при запросе
	}
	return user == false, nil // true, если пользователь не найден
}

//	func (s *UserService) CreateUser(req user.CreateUserRequest) (*user.User, error) {
//		newUser := &user.User{
//			Email:    req.Email,
//			Password: req.Password, // Здесь можно добавить хэширование
//		}
//
//		if err := s.Repo.Create(newUser); err != nil {
//			return nil, err
//		}
//		return newUser, nil
//	}
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hash), nil
}
func ComparePassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
func (s *UserService) Login(ctx context.Context, email, password string) (*user.UserResponse, error) {
	foundUser, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	if foundUser == nil {
		return nil, fmt.Errorf("user not found")
	}

	if !ComparePassword(foundUser.Password, password) {
		return nil, fmt.Errorf("invalid password")
	}

	// Возвращаем успешный ответ
	return &user.UserResponse{
		Name:  foundUser.Name,
		Email: foundUser.Email,
	}, nil
}
