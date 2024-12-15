package booking

import (
	"API/internal/lib/api/response"
	"API/internal/services/booking"
	bookingrepository "API/repositories/eventRepository"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type Handler struct {
	Logger   *slog.Logger
	Service  *booking.Service
	EventRep *bookingrepository.EventRep
}

// @Summary Получить мероприятие
// @Description Получить мероприятие по ID
// @Tags events
// @Success 200 {string} string "event"
// @Router / [get]
func (h *Handler) GetEventByID(w http.ResponseWriter, r *http.Request) {
	eventIDParam := chi.URLParam(r, "event_id")
	eventID, err := strconv.Atoi(eventIDParam)
	if err != nil {
		h.Logger.Error("invalid event_id format")
		render.JSON(w, r, response.Error("invalid event_id format"))
		return
	}
	event, err := h.EventRep.GetEvent(eventID)
	if err != nil {
		h.Logger.Error("failed to fetch event", slog.Int("event_id", eventID), slog.Any("error", err))
		render.JSON(w, r, response.Error("invalid fetch"))
		return
	}
	response.JSON(w, http.StatusOK, event)
}

func (h *Handler) AllEvents(w http.ResponseWriter, r *http.Request) {
	eventsList, err := h.EventRep.GetAllEvents()
	if err != nil {
		h.Logger.Error("failed to fetch list of events", err)
		render.JSON(w, r, response.Error("failed to fetch list of events"))
		return
	}
	response.JSON(w, http.StatusOK, eventsList)
}
