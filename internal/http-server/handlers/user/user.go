package userHandler

import (
	"API/internal/models/user"
	"API/internal/services/userService"
	"API/repositories/userRepository"
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
	"log/slog"
	"net/http"
	"os"
	"strconv"
)

type UserHandler struct {
	repo    *userRepository.UserRepository
	logger  *slog.Logger
	service *userService.UserService
}

// NewUserHandler создает новый экземпляр UserHandler
func NewUserHandler(repo *userRepository.UserRepository) *UserHandler {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	return &UserHandler{repo: repo, logger: logger}
}

// CreateUserHandler создает нового пользователя.
// @Summary      Create a new user
// @Description  Creates a new user with the provided information.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        user  body  user.CreateUserRequest  true  "User Data"
// @Success      200   {object}  user.UserResponse
// @Failure      400   {object}  map[string]string   "Invalid request body"
// @Failure      500   {object}  map[string]string   "Internal server error"
// @Router       /users [post]
func (h *UserHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	validate := validator.New()
	var req user.CreateUserRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		h.logger.Error("Invalid JSON input", "error", err)
		render.JSON(w, r, map[string]string{"error": "invalid request body"})
		return
	}
	if err := validate.Struct(req); err != nil {
		http.Error(w, "e-mail недействительный", http.StatusBadRequest)
		return
	}

	// Проверяем уникальность email
	isUnique, err := h.repo.IsEmailUnique(ctx, req.Email)
	if err != nil {
		h.logger.Error("Failed to check email uniqueness", "error", err)
		render.JSON(w, r, map[string]string{"error": "failed to validate email uniqueness"})
		return
	}
	if !isUnique {
		h.logger.Warn("Email is already in use", "email", req.Email)
		render.JSON(w, r, map[string]string{"error": "email is already in use"})
		return
	}

	// Создаем пользователя
	newUser, err := h.repo.NewUser(ctx, req)
	if err != nil {
		h.logger.Error("Failed to create user", "error", err)
		render.JSON(w, r, map[string]string{"error": "failed to create user"})
		return
	}

	h.logger.Info("User created successfully", "user_id", newUser.ID)
	render.JSON(w, r, map[string]interface{}{
		"user": map[string]interface{}{
			"id":    newUser.ID,
			"email": newUser.Email,
			"name":  newUser.Name,
		},
		"message": "User created successfully",
	})
}

// GetUserByIDHandler возвращает информацию о пользователе по его ID.
// @Summary      Get user by ID
// @Description  Retrieves information about a user by their ID.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id    path      int               true  "User ID"
// @Success      200   {object}  user.UserResponse
// @Failure      404   {object}  map[string]string "User not found"
// @Failure      500   {object}  map[string]string "Internal server error"
// @Router       /users/{id} [get]
func (h *UserHandler) GetUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Получаем ID из параметров запроса
	userID, err := parseUserID(r)
	if err != nil {
		h.logger.Warn("Invalid user ID", "error", err)
		render.JSON(w, r, map[string]string{"error": "invalid user ID"})
		return
	}

	// Ищем пользователя по ID
	user, err := h.repo.GetUserByID(ctx, userID)
	if err != nil {
		h.logger.Error("Failed to get user by ID", "error", err)
		render.JSON(w, r, map[string]string{"error": "user not found"})
		return
	}

	h.logger.Info("User retrieved successfully", "user_id", user.ID)
	render.JSON(w, r, map[string]interface{}{
		"user": map[string]interface{}{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
		},
	})
}

//// UpdateUserHandler обновляет данные пользователя.
//// @Summary      Update user
//// @Description  Updates information about an existing user.
//// @Tags         Users
//// @Accept       json
//// @Produce      json
//// @Param        id    path      int                   true  "User ID"
//// @Param        user  body      user.UpdateUserRequest  true  "Updated User Data"
//// @Success      200   {object}  user.UserResponse
//// @Failure      400   {object}  map[string]string     "Invalid request body"
//// @Failure      404   {object}  map[string]string     "User not found"
//// @Failure      500   {object}  map[string]string     "Internal server error"
//// @Router       /users/{id} [put]
//func (h *Handler) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
//	var req user.UpdateUserRequest
//	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
//		http.Error(w, "Invalid request body", http.StatusBadRequest)
//		return
//	}
//
//	id := chi.URLParam(r, "id")
//	updatedUser, err := h.repo.UpdateUser(r.Context(), id, req)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusNotFound)
//		return
//	}
//
//	w.WriteHeader(http.StatusOK)
//	json.NewEncoder(w).Encode(updatedUser)
//}

// DeleteUserHandler удаляет пользователя по его ID.
// @Summary      Delete user by ID
// @Description  Deletes a user by their ID.
// @Tags         Users
// @Param        id  path      int               true  "User ID"
// @Success      204  {object}  nil               "User deleted successfully"
// @Failure      404  {object}  map[string]string "User not found"
// @Failure      500  {object}  map[string]string "Internal server error"
// @Router       /users/{id} [delete]
func (h *UserHandler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Получаем ID из параметров запроса
	userID, err := parseUserID(r)
	if err != nil {
		h.logger.Warn("Invalid user ID", "error", err)
		render.JSON(w, r, map[string]string{"error": "invalid user ID"})
		return
	}

	// Удаляем пользователя
	if err := h.repo.DeleteUser(ctx, userID); err != nil {
		h.logger.Error("Failed to delete user", "error", err)
		render.JSON(w, r, map[string]string{"error": "failed to delete user"})
		return
	}

	h.logger.Info("User deleted successfully", "user_id", userID)
	render.JSON(w, r, map[string]string{
		"message": "User deleted successfully",
	})
}

// LoginHandler авторизует пользователя по email и паролю.
//
// @Summary Авторизация пользователя
// @Description Проверяет учетные данные пользователя и возвращает информацию о нем.
// @Tags Users
// @Accept json
// @Produce json
// @Param request body user.LoginUser true "Данные для авторизации"
// @Success 200 {object} user.LoginUser "Успешная авторизация"
// @Router /users/login [post]
func (h *UserHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req user.LoginUser
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		h.logger.Error("Invalid request body", "error", err)
		render.JSON(w, r, map[string]string{"error": "invalid request body"})
		return
	}

	userResp, err := h.service.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		h.logger.Error("Failed to login", "error", err)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	// Отправляем успешный ответ
	render.JSON(w, r, userResp)
}

// Вспомогательная функция для извлечения userID из параметров запроса
func parseUserID(r *http.Request) (uint, error) {
	userIDStr := chi.URLParam(r, "user_id")
	if userIDStr == "" {
		return 0, errors.New("user_id not provided in the URL")
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return 0, fmt.Errorf("invalid user_id format: %w", err)
	}

	if userID <= 0 {
		return 0, errors.New("user_id must be a positive integer")
	}

	return uint(userID), nil
}
