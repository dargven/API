package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims представляет данные, хранящиеся в JWT токене
type Claims struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// JWTManager управляет созданием и валидацией JWT токенов
type JWTManager struct {
	secret   []byte
	tokenTTL time.Duration
}

// NewJWTManager создает новый экземпляр JWTManager
func NewJWTManager(secret string, tokenTTL time.Duration) *JWTManager {
	return &JWTManager{
		secret:   []byte(secret),
		tokenTTL: tokenTTL,
	}
}

// GenerateToken создает новый JWT токен для пользователя
func (m *JWTManager) GenerateToken(userID int64, email string) (string, error) {
	const op = "auth.JWTManager.GenerateToken"

	now := time.Now()

	claims := Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(m.tokenTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "api-curc",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(m.secret)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return tokenString, nil
}

// ValidateToken проверяет JWT токен и возвращает claims
func (m *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
	const op = "auth.JWTManager.ValidateToken"

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Проверяем алгоритм подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.secret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("%s: invalid token", op)
	}

	return claims, nil
}
