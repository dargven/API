package bookings

import (
	authMiddleware "API/internal/http-server/middleware/auth"
	resp "API/internal/lib/api/response"
	"API/internal/lib/logger/sl"
	"API/internal/models"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"golang.org/x/exp/slog"
)

// BookingsLister интерфейс для получения списка бронирований
type BookingsLister interface {
	GetBookingsByUserID(userID int64) ([]*models.BookingWithEvent, error)
}

// ListBookingsResponse ответ со списком бронирований
type ListBookingsResponse struct {
	resp.Response
	Bookings []models.BookingResponse `json:"bookings"`
}

// NewList возвращает хендлер для получения списка бронирований
// @Summary Мои билеты
// @Description Возвращает список всех бронирований пользователя
// @Tags bookings
// @Security BearerAuth
// @Produce json
// @Success 200 {object} ListBookingsResponse
// @Failure 401 {object} resp.Response
// @Failure 500 {object} resp.Response
// @Router /bookings [get]
func NewList(log *slog.Logger, lister BookingsLister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.bookings.List"

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

		bookings, err := lister.GetBookingsByUserID(userID)
		if err != nil {
			log.Error("failed to get bookings", sl.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to get bookings"))
			return
		}

		response := make([]models.BookingResponse, 0, len(bookings))
		for _, b := range bookings {
			response = append(response, b.ToResponse())
		}

		render.JSON(w, r, ListBookingsResponse{
			Response: resp.OK(),
			Bookings: response,
		})
	}
}
