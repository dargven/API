package events

import (
	storage "API/internal/Storage"
	resp "API/internal/lib/api/response"
	"API/internal/lib/logger/sl"
	"API/internal/models"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"golang.org/x/exp/slog"
)

// EventGetter интерфейс для получения мероприятий
type EventGetter interface {
	GetEventByID(id int64) (*models.Event, error)
	GetAllEvents(limit, offset int) ([]*models.Event, error)
}

// GetByIDResponse структура ответа при получении мероприятия по ID
// @Description Ответ при получении мероприятия по ID
type GetByIDResponse struct {
	resp.Response
	Event models.EventResponse `json:"event"`
}

// GetAllResponse структура ответа при получении всех мероприятий
// @Description Ответ при получении списка мероприятий
type GetAllResponse struct {
	resp.Response
	Events []models.EventResponse `json:"events"`
	Total  int                    `json:"total" example:"10"`
}

// NewGetByID создает хендлер для получения мероприятия по ID
// @Summary Получение мероприятия по ID
// @Description Возвращает мероприятие по его ID. Требуется JWT авторизация.
// @Tags events
// @Produce json
// @Param Authorization header string true "Bearer JWT токен"
// @Param id path int true "ID мероприятия"
// @Success 200 {object} GetByIDResponse
// @Failure 400 {object} resp.Response
// @Failure 401 {object} resp.Response
// @Failure 404 {object} resp.Response
// @Failure 500 {object} resp.Response
// @Router /events/{id} [get]
func NewGetByID(log *slog.Logger, eventGetter EventGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.events.GetByID"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		idStr := chi.URLParam(r, "id")
		if idStr == "" {
			log.Error("id parameter is empty")
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("id parameter is required"))
			return
		}

		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Error("invalid id format", sl.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid id format"))
			return
		}

		event, err := eventGetter.GetEventByID(id)
		if errors.Is(err, storage.ErrEventNotFound) {
			log.Info("event not found", slog.Int64("id", id))
			w.WriteHeader(http.StatusNotFound)
			render.JSON(w, r, resp.Error("event not found"))
			return
		}
		if err != nil {
			log.Error("failed to get event", sl.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("internal error"))
			return
		}

		log.Info("event retrieved", slog.Int64("event_id", event.ID))

		render.JSON(w, r, GetByIDResponse{
			Response: resp.OK(),
			Event:    event.ToResponse(),
		})
	}
}

// NewGetAll создает хендлер для получения всех мероприятий
// @Summary Получение списка мероприятий
// @Description Возвращает список всех мероприятий с пагинацией. Требуется JWT авторизация.
// @Tags events
// @Produce json
// @Param Authorization header string true "Bearer JWT токен"
// @Param limit query int false "Количество записей (по умолчанию 20, максимум 100)"
// @Param offset query int false "Смещение (по умолчанию 0)"
// @Success 200 {object} GetAllResponse
// @Failure 401 {object} resp.Response
// @Failure 500 {object} resp.Response
// @Router /events [get]
func NewGetAll(log *slog.Logger, eventGetter EventGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.events.GetAll"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		// Пагинация
		limitStr := r.URL.Query().Get("limit")
		offsetStr := r.URL.Query().Get("offset")

		limit := 20 // default
		offset := 0

		if limitStr != "" {
			if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
				limit = l
			}
		}

		if offsetStr != "" {
			if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
				offset = o
			}
		}

		events, err := eventGetter.GetAllEvents(limit, offset)
		if err != nil {
			log.Error("failed to get events", sl.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("internal error"))
			return
		}

		// Конвертируем в response формат
		eventResponses := make([]models.EventResponse, 0, len(events))
		for _, event := range events {
			eventResponses = append(eventResponses, event.ToResponse())
		}

		log.Info("events retrieved", slog.Int("count", len(events)))

		render.JSON(w, r, GetAllResponse{
			Response: resp.OK(),
			Events:   eventResponses,
			Total:    len(eventResponses),
		})
	}
}
