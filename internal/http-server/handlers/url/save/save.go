package save

import (
	// для краткости даем короткий алиас пакету
	storage "API/internal/Storage"
	resp "API/internal/lib/api/response"
	"API/internal/lib/logger/sl"
	"API/internal/lib/random"
	"errors"
	"io"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
	"golang.org/x/exp/slog"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}
type URLSaver interface {
	SaveURL(URL, alias string) (int64, error)
}

const aliasLenght = 6

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		// к логеру добавляем поля op & request_id

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			// ошибка вылезет если запрос будет с пустым телом
			log.Error("request body is empty")

			render.JSON(w, r, resp.Error("empty request"))

			return

		}
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode body"))

			return

		}

		log.Info("request body decoded", slog.Any("req", req))

		// Создаем объект валидатора
		// и передаем в него структуру, которую нужно провалидировать
		if err := validator.New().Struct(req); err != nil {
			validateError := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, resp.Error(validateError.Error()))
		}

		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(aliasLenght)
		}

		id, err := urlSaver.SaveURL(req.URL, alias)
		if errors.Is(err, storage.ErrURLExists) {
			// Отдельно обрабатываем ситуацию,
			// когда запись с таким Alias уже существует
			log.Info("url already exist", slog.String("url", req.URL))

			render.JSON(w, r, resp.Error("this url is already exist"))

			return
		}
		if err != nil {
			log.Error("failed to add url", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to add url"))

			return
		}
		log.Info("url added", slog.Int64("id", id))

		responseOK(w, r, alias)
	}

}

func responseOK(w http.ResponseWriter, r *http.Request, alias string) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Alias:    alias,
	})
}
