package events

import (
	authMiddleware "API/internal/http-server/middleware/auth"
	resp "API/internal/lib/api/response"
	"API/internal/lib/logger/sl"
	"API/internal/models"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"golang.org/x/exp/slog"
)

// EventCreator интерфейс для создания мероприятий
type EventCreator interface {
	CreateEvent(event *models.Event) (*models.Event, error)
}

// CreateRequest структура запроса на создание мероприятия
// @Description Запрос на создание мероприятия
type CreateRequest struct {
	Title       string `json:"title" validate:"required,min=3,max=200" example:"Конференция Go"`
	Description string `json:"description" validate:"max=2000" example:"Ежегодная конференция разработчиков Go"`
	Location    string `json:"location" validate:"required,max=500" example:"Москва, ул. Примерная 1"`
	StartTime   string `json:"start_time" validate:"required" example:"2024-06-15T10:00:00Z"`
	EndTime     string `json:"end_time" validate:"required" example:"2024-06-15T18:00:00Z"`
	MaxSlots    int    `json:"max_slots" validate:"required,min=1" example:"100"`
}

// CreateResponse структура ответа при создании мероприятия
// @Description Ответ при создании мероприятия
type CreateResponse struct {
	resp.Response
	Event models.EventResponse `json:"event"`
}

// NewCreate создает хендлер для создания мероприятия
// @Summary Создание мероприятия
// @Description Создает новое мероприятие. Требуется JWT авторизация.
// @Tags events
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT токен"
// @Param request body CreateRequest true "Данные мероприятия"
// @Success 201 {object} CreateResponse
// @Failure 400 {object} resp.Response
// @Failure 401 {object} resp.Response
// @Failure 500 {object} resp.Response
// @Router /events [post]
func NewCreate(log *slog.Logger, eventCreator EventCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.events.Create"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		userID, ok := authMiddleware.GetUserIDFromContext(r.Context())
		if !ok {
			log.Error("user_id not found in context")
			w.WriteHeader(http.StatusUnauthorized)
			render.JSON(w, r, resp.Error("unauthorized"))
			return
		}

		var req CreateRequest
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid request body"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		// Валидация
		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("validation failed", sl.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.ValidationError(validateErr))
			return
		}

		// Парсинг времени
		startTime, err := time.Parse(time.RFC3339, req.StartTime)
		if err != nil {
			log.Error("invalid start_time format", sl.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid start_time format, use RFC3339 (e.g., 2024-01-15T10:00:00Z)"))
			return
		}

		endTime, err := time.Parse(time.RFC3339, req.EndTime)
		if err != nil {
			log.Error("invalid end_time format", sl.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid end_time format, use RFC3339 (e.g., 2024-01-15T12:00:00Z)"))
			return
		}

		// Проверка что end_time > start_time
		if !endTime.After(startTime) {
			log.Error("end_time must be after start_time")
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("end_time must be after start_time"))
			return
		}

		// Проверка что start_time в будущем
		if startTime.Before(time.Now()) {
			log.Error("start_time must be in the future")
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("start_time must be in the future"))
			return
		}

		event := &models.Event{
			Title:       req.Title,
			Description: req.Description,
			Location:    req.Location,
			StartTime:   startTime,
			EndTime:     endTime,
			CreatorID:   userID,
			MaxSlots:    req.MaxSlots,
		}

		createdEvent, err := eventCreator.CreateEvent(event)
		if err != nil {
			log.Error("failed to create event", sl.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to create event"))
			return
		}

		log.Info("event created", slog.Int64("event_id", createdEvent.ID))

		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, CreateResponse{
			Response: resp.OK(),
			Event:    createdEvent.ToResponse(),
		})
	}
}
