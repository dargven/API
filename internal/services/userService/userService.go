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

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hash), nil
}
func (s *UserService) RegisterUser(ctx context.Context, req user.CreateUserRequest) (*user.User, error) {
	// Проверяем уникальность email
	isUnique, err := s.repo.IsEmailUnique(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check email uniqueness: %w", err)
	}
	if !isUnique {
		return nil, errors.New("email is already in use")
	}

	// Хэшируем пароль
	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	req.Password = hashedPassword

	// Создаем пользователя в базе данных
	newUser, err := s.repo.NewUser(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return newUser, nil
}

func ComparePassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func (s *UserService) Login(ctx context.Context, email, password string) (*user.User, error) {
	// Проверяем, есть ли пользователь с данным email
	foundUser, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || foundUser == nil {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	// Проверяем пароль
	if !ComparePassword(foundUser.Password, password) {
		return nil, errors.New("invalid password")
	}

	// Успешная авторизация, возвращаем данные пользователя
	return foundUser, nil
}
