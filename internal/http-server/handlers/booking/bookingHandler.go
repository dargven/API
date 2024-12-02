package booking

import (
	"API/internal/lib/api/response"
	"API/internal/services/booking"
	bookingrepository "API/repositories/bookingRepository"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"strconv"
)

type Handler struct {
	Logger    *slog.Logger
	Service   *booking.Service
	EventRepo *bookingrepository.BookingRep
}

func (h *Handler) GetEventByID(w http.ResponseWriter, r *http.Request) {
	eventIDParam := chi.URLParam(r, "event_id")
	eventID, err := strconv.Atoi(eventIDParam)
	if err != nil {
		h.Logger.Error("invalid event_id format")
		render.JSON(w, r, response.Error("invalid event_id format"))
		return
	}
	event, err := h.EventRepo.FetchEventByID(eventID)
	if err != nil {
		h.Logger.Error("failed to fetch event", slog.Int("event_id", eventID), slog.Any("error", err))
		render.JSON(w, r, response.Error("invalid fetch"))
		return
	}
	response.JSON(w, http.StatusOK, event)
}
