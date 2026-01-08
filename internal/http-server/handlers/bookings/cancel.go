package bookings

import (
	storage "API/internal/Storage"
	authMiddleware "API/internal/http-server/middleware/auth"
	resp "API/internal/lib/api/response"
	"API/internal/lib/logger/sl"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"golang.org/x/exp/slog"
)

// BookingCanceller интерфейс для отмены бронирования
type BookingCanceller interface {
	CancelBooking(bookingID, userID int64) error
}

// NewCancel возвращает хендлер для отмены бронирования
// @Summary Отменить бронь
// @Description Отменяет бронирование и возвращает деньги на баланс
// @Tags bookings
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID бронирования"
// @Success 200 {object} resp.Response
// @Failure 400 {object} resp.Response
// @Failure 401 {object} resp.Response
// @Failure 404 {object} resp.Response
// @Failure 500 {object} resp.Response
// @Router /bookings/{id} [delete]
func NewCancel(log *slog.Logger, canceller BookingCanceller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.bookings.Cancel"

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

		bookingIDStr := chi.URLParam(r, "id")
		bookingID, err := strconv.ParseInt(bookingIDStr, 10, 64)
		if err != nil {
			log.Error("invalid booking id", sl.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid booking id"))
			return
		}

		err = canceller.CancelBooking(bookingID, userID)
		if err != nil {
			if errors.Is(err, storage.ErrBookingNotFound) {
				log.Error("booking not found", slog.Int64("booking_id", bookingID))
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, resp.Error("booking not found"))
				return
			}

			log.Error("failed to cancel booking", sl.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to cancel booking"))
			return
		}

		log.Info("booking cancelled", slog.Int64("booking_id", bookingID))
		render.JSON(w, r, resp.OK())
	}
}
