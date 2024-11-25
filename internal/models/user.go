package models

import (
	resp "API/internal/lib/api/response"
	"API/internal/lib/logger/sl"
	"context"
	"database/sql"
	"errors"
	"io"
	"net/http"
	"os"
	"regexp"

	"log/slog"

	"API/internal/Storage/postrgeSQL"

	"github.com/go-chi/render"
	"github.com/go-playground/validator"
)

// здесь же описан сам пользак и функции по его добавлению с проверками и отправкой всего в бд

type UserHandler struct {
	DB *postrgeSQL.Database
}

type User struct {
	ID       int64
	Email    string `validate:"required,email"`
	Name     string `validate:"required"`
	Password string `validate:"required"`
}

func (h *UserHandler) NewUser(w http.ResponseWriter, r *http.Request) {
	var user User
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	err := render.DecodeJSON(r.Body, &user)
	if errors.Is(err, io.EOF) {
		//Если пусто
		logger.Error("body is empty")

		render.JSON(w, r, resp.Error("body is empty"))

		return
	}

	if err != nil {
		//ошибка при попытке декода
		logger.Error("failed to decode body", sl.Err(err))

		render.JSON(w, r, resp.Error("failed to decode"))

		return
	}

	if err := validator.New().Struct(user); err != nil {
		validateError := err.(validator.ValidationErrors)

		logger.Error(validateError.Error(), sl.Err(err))

		render.JSON(w, r, resp.Error("invalid request"))
	}

	if user.Email == "" || user.Name == "" || user.Password == "" {
		logger.Error("all fields must be filled")

		render.JSON(w, r, resp.Error("fields are not filled"))

		return
	}

	if !isEmailValid(user.Email) {
		logger.Error("email does not match the format")

		render.JSON(w, r, resp.Error("email does not match the format"))

		return
	}

	if h.isUnique(user.Email) {
		logger.Error("email already exist")

		render.JSON(w, r, resp.Error("email already exist"))

		return
	}

	id, err := h.DB.AddUser(user.Email, user.Name, user.Password)
	if err != nil {
		logger.Error("failed to create user")

		render.JSON(w, r, resp.Error("failed to create user"))

		return
	}

	render.JSON(w, r, id)

}

func (h *UserHandler) isUnique(email string) bool {
	query := `SELECT id FROM users WHERE email = $1`

	var id int64

	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	err := h.DB.Pool.QueryRow(context.Background(), query, email).Scan(&id)
	if err != nil {
		logger.Error("Error checking email uniqueness:", err)
		return false
	}
	if err == sql.ErrNoRows {
		return false
	}

	return true
}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(e)
}
