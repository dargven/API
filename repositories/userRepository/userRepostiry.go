package userRepository

import (
	"API/internal/Storage/postrgeSQL"
	resp "API/internal/lib/api/response"
	"API/internal/models/user"
	"errors"
	"github.com/go-chi/render"
	"io"
	"log/slog"
	"net/http"
	"os"
)

type userRepository struct {
	db *postrgeSQL.Database
}

func (h *userRepository) NewUser(w http.ResponseWriter, r *http.Request) {
	var userRequest user.CreateUserRequest
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	err := render.DecodeJSON(r.Body, &userRequest)
	if errors.Is(err, io.EOF) {
		//Если пусто
		logger.Error("body is empty")

		render.JSON(w, r, resp.Error("body is empty"))

		return
	}
