package test

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
)

func HelloHandler(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.hello"

		log = log.With(slog.String("op", op))
		log.Info("handler called")

		response := map[string]string{"message": "working"}
		render.JSON(w, r, response)
	}
}
