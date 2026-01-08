package bookings

import (
	storage "API/internal/Storage"
	authMiddleware "API/internal/http-server/middleware/auth"
	resp "API/internal/lib/api/response"
	"API/internal/lib/logger/sl"
	"API/internal/models"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	chimiddleware "github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"golang.org/x/exp/slog"
)

// BookingCreator интерфейс для создания бронирования
type BookingCreator interface {
	CreateBooking(userID, eventID int64, quantity int) (*models.Booking, error)
}

// CreateBookingRequest запрос на создание бронирования
type CreateBookingRequest struct {
	Quantity int `json:"quantity" validate:"required,min=1,max=10"`
}

// CreateBookingResponse ответ с данными бронирования
type CreateBookingResponse struct {
	resp.Response
	Booking models.Booking `json:"booking"`
}

// NewCreate возвращает хендлер для создания бронирования
// @Summary Забронировать билет
// @Description Бронирует билет на мероприятие
// @Tags bookings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID мероприятия"
// @Param request body CreateBookingRequest true "Количество билетов"
// @Success 201 {object} CreateBookingResponse
// @Failure 400 {object} resp.Response
// @Failure 401 {object} resp.Response
// @Failure 404 {object} resp.Response
// @Failure 409 {object} resp.Response "Бронирование уже существует"
// @Failure 422 {object} resp.Response "Недостаточно билетов или баланса"
// @Failure 500 {object} resp.Response
// @Router /events/{id}/book [post]
func NewCreate(log *slog.Logger, creator BookingCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.bookings.Create"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", chimiddleware.GetReqID(r.Context())),
		)

		userID, ok := authMiddleware.GetUserIDFromContext(r.Context())
		if !ok {
			log.Error("user_id not found in context")
			w.WriteHeader(http.StatusUnauthorized)
			render.JSON(w, r, resp.Error("unauthorized"))
			return
		}

		eventIDStr := chi.URLParam(r, "id")
		eventID, err := strconv.ParseInt(eventIDStr, 10, 64)
		if err != nil {
			log.Error("invalid event id", sl.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid event id"))
			return
		}

		var req CreateBookingRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode request", sl.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid request body"))
			return
		}

		validate := validator.New()
		if err := validate.Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("validation failed", sl.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.ValidationError(validateErr))
			return
		}

		booking, err := creator.CreateBooking(userID, eventID, req.Quantity)
		if err != nil {
			if errors.Is(err, storage.ErrEventNotFound) {
				log.Error("event not found", slog.Int64("event_id", eventID))
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, resp.Error("event not found"))
				return
			}
			if errors.Is(err, storage.ErrNoTickets) {
				log.Error("no available tickets", slog.Int64("event_id", eventID))
				w.WriteHeader(http.StatusUnprocessableEntity)
				render.JSON(w, r, resp.Error("not enough available tickets"))
				return
			}
			if errors.Is(err, storage.ErrInsufficientBalance) {
				log.Error("insufficient balance", slog.Int64("user_id", userID))
				w.WriteHeader(http.StatusUnprocessableEntity)
				render.JSON(w, r, resp.Error("insufficient balance"))
				return
			}
			if errors.Is(err, storage.ErrBookingExists) {
				log.Error("booking already exists", slog.Int64("user_id", userID), slog.Int64("event_id", eventID))
				w.WriteHeader(http.StatusConflict)
				render.JSON(w, r, resp.Error("you already have a booking for this event"))
				return
			}

			log.Error("failed to create booking", sl.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to create booking"))
			return
		}

		log.Info("booking created",
			slog.Int64("booking_id", booking.ID),
			slog.String("booking_code", booking.BookingCode),
		)

		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, CreateBookingResponse{
			Response: resp.OK(),
			Booking:  *booking,
		})
	}
}
