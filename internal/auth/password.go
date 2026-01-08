package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const (
	// bcrypt cost - баланс между безопасностью и производительностью
	bcryptCost = 10
)

// HashPassword хеширует пароль с использованием bcrypt
func HashPassword(password string) (string, error) {
	const op = "auth.HashPassword"

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return string(bytes), nil
}

// CheckPassword сравнивает пароль с хешем (constant-time comparison)
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
