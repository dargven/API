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

// UserGetter интерфейс для получения пользователя
type UserGetter interface {
	GetUserByEmail(email string) (*models.User, error)
}

// LoginRequest запрос на авторизацию
// @Description Запрос на авторизацию
type LoginRequest struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"securePassword123"`
}

// LoginResponse ответ при успешной авторизации
// @Description Ответ при успешной авторизации
type LoginResponse struct {
	resp.Response
	Token string `json:"token,omitempty" example:"eyJhbGciOiJIUzI1NiIs..."`
}

// NewLogin создает хендлер авторизации
// @Summary Авторизация пользователя
// @Description Авторизует пользователя и возвращает JWT токен
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Данные для авторизации"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} resp.Response
// @Failure 401 {object} resp.Response
// @Failure 500 {object} resp.Response
// @Router /auth/login [post]
func NewLogin(log *slog.Logger, userGetter UserGetter, jwtManager *auth.JWTManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.auth.Login"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req LoginRequest
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
		if req.Email == "" || req.Password == "" {
			log.Error("missing required fields")
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("email and password are required"))
			return
		}

		if !models.IsEmailValid(req.Email) {
			log.Error("invalid email format")
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid email format"))
			return
		}

		// Получаем пользователя по email
		user, err := userGetter.GetUserByEmail(req.Email)
		if err != nil {
			if errors.Is(err, storage.ErrUserNotFound) {
				log.Info("user not found", slog.String("email", req.Email))
				w.WriteHeader(http.StatusUnauthorized)
				render.JSON(w, r, resp.Error("invalid email or password"))
				return
			}
			log.Error("failed to get user", sl.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("internal error"))
			return
		}

		// Проверяем пароль
		if !auth.CheckPassword(req.Password, user.PasswordHash) {
			log.Info("invalid password", slog.String("email", req.Email))
			w.WriteHeader(http.StatusUnauthorized)
			render.JSON(w, r, resp.Error("invalid email or password"))
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

		log.Info("user logged in successfully", slog.Int64("user_id", user.ID))

		render.JSON(w, r, LoginResponse{
			Response: resp.OK(),
			Token:    token,
		})
	}
}
