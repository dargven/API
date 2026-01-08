package auth

import (
	"API/internal/Storage"
	"API/internal/auth"
	resp "API/internal/lib/api/response"
	"API/internal/lib/logger/sl"
	"API/internal/models"
	"errors"
	"io"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"golang.org/x/exp/slog"
)

// UserCreator интерфейс для создания пользователей
type UserCreator interface {
	CreateUser(email, name, passwordHash string) (*models.User, error)
	EmailExists(email string) (bool, error)
}

// RegisterRequest запрос на регистрацию
// @Description Запрос на регистрацию пользователя
type RegisterRequest struct {
	Email    string `json:"email" example:"user@example.com"`
	Name     string `json:"name" example:"John Doe"`
	Password string `json:"password" example:"securePassword123"`
}

// RegisterResponse ответ при успешной регистрации
// @Description Ответ при успешной регистрации
type RegisterResponse struct {
	resp.Response
	User  models.UserResponse `json:"user,omitempty"`
	Token string              `json:"token,omitempty" example:"eyJhbGciOiJIUzI1NiIs..."`
}

// NewRegister создает хендлер регистрации
// @Summary Регистрация пользователя
// @Description Создает нового пользователя и возвращает JWT токен
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Данные для регистрации"
// @Success 201 {object} RegisterResponse
// @Failure 400 {object} resp.Response
// @Failure 409 {object} resp.Response
// @Failure 500 {object} resp.Response
// @Router /auth/register [post]
func NewRegister(log *slog.Logger, userCreator UserCreator, jwtManager *auth.JWTManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.auth.Register"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req RegisterRequest
		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("empty request body"))
			return
		}
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}

		// Валидация полей
		if req.Email == "" || req.Name == "" || req.Password == "" {
			log.Error("missing required fields")
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("email, name and password are required"))
			return
		}

		if !models.IsEmailValid(req.Email) {
			log.Error("invalid email format", slog.String("email", req.Email))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid email format"))
			return
		}

		if !models.IsPasswordValid(req.Password) {
			log.Error("password too short")
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("password must be at least 8 characters"))
			return
		}

		// Проверяем, не занят ли email
		exists, err := userCreator.EmailExists(req.Email)
		if err != nil {
			log.Error("failed to check email existence", sl.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("internal error"))
			return
		}
		if exists {
			log.Info("email already exists", slog.String("email", req.Email))
			w.WriteHeader(http.StatusConflict)
			render.JSON(w, r, resp.Error("user with this email already exists"))
			return
		}

		// Хешируем пароль
		passwordHash, err := auth.HashPassword(req.Password)
		if err != nil {
			log.Error("failed to hash password", sl.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("internal error"))
			return
		}

		// Создаем пользователя
		user, err := userCreator.CreateUser(req.Email, req.Name, passwordHash)
		if err != nil {
			if errors.Is(err, storage.ErrUserExists) {
				log.Info("user already exists", slog.String("email", req.Email))
				w.WriteHeader(http.StatusConflict)
				render.JSON(w, r, resp.Error("user with this email already exists"))
				return
			}
			log.Error("failed to create user", sl.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to create user"))
			return
		}

		// Генерируем JWT токен
		token, err := jwtManager.GenerateToken(user.ID, user.Email)
		if err != nil {
			log.Error("failed to generate token", sl.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to generate token"))
			return
		}

		log.Info("user registered successfully", slog.Int64("user_id", user.ID))

		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, RegisterResponse{
			Response: resp.OK(),
			User:     user.ToResponse(),
			Token:    token,
		})
	}
}
