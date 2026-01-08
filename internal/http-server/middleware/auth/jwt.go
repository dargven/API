package auth

import (
	"API/internal/auth"
	resp "API/internal/lib/api/response"
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/render"
	"golang.org/x/exp/slog"
)

// ContextKey тип для ключей контекста
type ContextKey string

const (
	// UserIDKey ключ для user_id в контексте
	UserIDKey ContextKey = "user_id"
	// UserEmailKey ключ для email в контексте
	UserEmailKey ContextKey = "user_email"
)

// JWTAuth создает middleware для проверки JWT токена
func JWTAuth(logger *slog.Logger, jwtManager *auth.JWTManager) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			const op = "middleware.auth.JWTAuth"

			log.Printf("[v0] JWTAuth middleware called for %s %s", r.Method, r.URL.Path)

			// Получаем токен из заголовка Authorization
			authHeader := r.Header.Get("Authorization")
			log.Printf("[v0] Authorization header: %q", authHeader)

			if authHeader == "" {
				logger.Info("missing authorization header", slog.String("op", op))
				w.WriteHeader(http.StatusUnauthorized)
				render.JSON(w, r, resp.Error("authorization header is required"))
				return
			}

			authHeader = strings.TrimSpace(authHeader)
			parts := strings.SplitN(authHeader, " ", 2)
			log.Printf("[v0] Header parts: %v (len=%d)", parts, len(parts))

			if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
				logger.Info("invalid authorization header format", slog.String("op", op))
				w.WriteHeader(http.StatusUnauthorized)
				render.JSON(w, r, resp.Error("invalid authorization header format"))
				return
			}

			tokenString := strings.TrimSpace(parts[1])
			log.Printf("[v0] Token string length: %d", len(tokenString))

			// Валидируем токен
			claims, err := jwtManager.ValidateToken(tokenString)
			if err != nil {
				log.Printf("[v0] Token validation error: %v", err)
				logger.Info("invalid token", slog.String("op", op), slog.String("error", err.Error()))
				w.WriteHeader(http.StatusUnauthorized)
				render.JSON(w, r, resp.Error("invalid or expired token"))
				return
			}

			log.Printf("[v0] Token valid! UserID: %d, Email: %s", claims.UserID, claims.Email)

			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, UserEmailKey, claims.Email)

			// Передаем запрос дальше с обновленным контекстом
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserIDFromContext извлекает user_id из контекста
func GetUserIDFromContext(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(UserIDKey).(int64)
	return userID, ok
}

// GetUserEmailFromContext извлекает email из контекста
func GetUserEmailFromContext(ctx context.Context) (string, bool) {
	email, ok := ctx.Value(UserEmailKey).(string)
	return email, ok
}
