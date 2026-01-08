package models

import (
	"regexp"
	"time"
)

// User представляет пользователя системы
type User struct {
	ID           int64     `json:"id"`
	Email        string    `json:"email" validate:"required,email"`
	Name         string    `json:"name" validate:"required"`
	PasswordHash string    `json:"-"` // не отдаем в JSON
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// UserResponse - DTO для ответа без чувствительных данных
type UserResponse struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// ToResponse конвертирует User в UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		CreatedAt: u.CreatedAt,
	}
}

// IsEmailValid проверяет формат email
func IsEmailValid(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// IsPasswordValid проверяет минимальные требования к паролю
func IsPasswordValid(password string) bool {
	return len(password) >= 8
}
