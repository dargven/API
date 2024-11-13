package models

import (
	resp "API/internal/lib/api/response"
	"API/internal/lib/logger/sl"
	"errors"
	"io"
	"net/http"
	"os"
	"regexp"

	"log/slog"

	"github.com/go-chi/render"
	"github.com/go-playground/validator"
)

// здесь же описан сам пользак и функции по его добавлению с проверками и отправкой всего в бд

type User struct {
	ID       int64
	Email    string `validate:"required,email"`
	Name     string `validate:"required"`
	Password string `validate:"required"`
}

func newUser(w http.ResponseWriter, r *http.Request) {
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

		logger.Error("invalid request", sl.Err(err))

		render.JSON(w, r, resp.Error(validateError.Error()))
	}

	// if user.Email == "" || user.Name == "" || user.Password == "" {
	// 	logger.Error("all fields must be filled")

	// 	render.JSON(w, r, resp.Error("fields are not filled"))

	// 	return
	// }

	if !isEmailValid(user.Email) {
		logger.Error("email does not match the format")

		render.JSON(w, r, resp.Error("email does not match the format"))

		return
	} else if isUnique(user.Email) { //уникальность
		logger.Error("email already exist")

		render.JSON(w, r, resp.Error("email already exist"))

		return
	} // по такому же принципу надо сделать проверку на совпадение

	//по желанию сложность пароля тоже тут надо реализовать

}

func isUnique(email string) bool {
	//проверка на уникальность будет реализованна когда поднимем бд
	return false
}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(e)
}
